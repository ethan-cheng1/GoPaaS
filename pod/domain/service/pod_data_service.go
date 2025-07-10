package service

import (
	"context"
	"errors"
	"git.imooc.com/coding-535/common"
	"git.imooc.com/coding-535/pod/domain/model"
	"git.imooc.com/coding-535/pod/domain/repository"
	"git.imooc.com/coding-535/pod/proto/pod"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

type IPodDataService interface {
	AddPod(*model.Pod) (int64, error)
	DeletePod(int64) error
	UpdatePod(*model.Pod) error
	FindPodByID(int64) (*model.Pod, error)
	FindAllPod() ([]model.Pod, error)
	CreateToK8s(*pod.PodInfo) error
	DeleteFromK8s(*model.Pod) error
	UpdateToK8s(*pod.PodInfo) error
}

func NewPodDataService(podRepository repository.IPodRepository, clientSet *kubernetes.Clientset) IPodDataService {
	return &PodDataService{
		PodRepository: podRepository,
		K8sClientSet:  clientSet,
		deployment:    &v1.Deployment{},
	}
}

type PodDataService struct {
	PodRepository repository.IPodRepository
	K8sClientSet  *kubernetes.Clientset
	deployment    *v1.Deployment
}

// Create pod in k8s
func (u *PodDataService) CreateToK8s(podInfo *pod.PodInfo) (err error) {
	u.SetDeployment(podInfo)
	if _, err = u.K8sClientSet.AppsV1().Deployments(podInfo.PodNamespace).Get(context.TODO(), podInfo.PodName, v12.GetOptions{}); err != nil {
		if _, err = u.K8sClientSet.AppsV1().Deployments(podInfo.PodNamespace).Create(context.TODO(), u.deployment, v12.CreateOptions{}); err != nil {
			common.Error(err)
			return err
		}
		common.Info("Created successfully")
		return nil
	} else {
		// Can write custom business logic here
		common.Error("Pod " + podInfo.PodName + " already exists")
		return errors.New("Pod " + podInfo.PodName + " already exists")
	}

}

// Update deployment and pod
func (u *PodDataService) UpdateToK8s(podInfo *pod.PodInfo) (err error) {
	u.SetDeployment(podInfo)
	if _, err = u.K8sClientSet.AppsV1().Deployments(podInfo.PodNamespace).Get(context.TODO(), podInfo.PodName, v12.GetOptions{}); err != nil {
		common.Error(err)
		return errors.New("Pod " + podInfo.PodName + " does not exist, please create first")
	} else {
		// If exists
		if _, err = u.K8sClientSet.AppsV1().Deployments(podInfo.PodNamespace).Update(context.TODO(), u.deployment, v12.UpdateOptions{}); err != nil {
			common.Error(err)
			return err
		}
		common.Info(podInfo.PodName + " updated successfully")
		return nil
	}

}

// Delete pod
func (u *PodDataService) DeleteFromK8s(pod *model.Pod) (err error) {
	if err = u.K8sClientSet.AppsV1().Deployments(pod.PodNamespace).Delete(context.TODO(), pod.PodName, v12.DeleteOptions{}); err != nil {
		common.Error(err)
		// Can write custom business logic here
		return err
	} else {
		if err := u.DeletePod(pod.ID); err != nil {
			common.Error(err)
			return err
		}
		common.Info("Deleted Pod ID: " + strconv.FormatInt(pod.ID, 10) + " successfully!")
	}
	return
}

func (u *PodDataService) SetDeployment(podInfo *pod.PodInfo) {
	deployment := &v1.Deployment{}
	deployment.TypeMeta = v12.TypeMeta{
		Kind:       "deployment",
		APIVersion: "v1",
	}
	deployment.ObjectMeta = v12.ObjectMeta{
		Name:      podInfo.PodName,
		Namespace: podInfo.PodNamespace,
		Labels: map[string]string{
			"app-name": podInfo.PodName,
			"author":   "Caplost",
		},
	}
	deployment.Name = podInfo.PodName
	deployment.Spec = v1.DeploymentSpec{
		// Number of replicas
		Replicas: &podInfo.PodReplicas,
		Selector: &v12.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": podInfo.PodName,
			},
			MatchExpressions: nil,
		},
		Template: v13.PodTemplateSpec{
			ObjectMeta: v12.ObjectMeta{
				Labels: map[string]string{
					"app-name": podInfo.PodName,
				},
			},
			Spec: v13.PodSpec{
				Containers: []v13.Container{
					{
						Name:            podInfo.PodName,
						Image:           podInfo.PodImage,
						Ports:           u.getContainerPort(podInfo),
						Env:             u.getEnv(podInfo),
						Resources:       u.getResources(podInfo),
						ImagePullPolicy: u.getImagePullPolicy(podInfo),
					},
				},
			},
		},
		Strategy:                v1.DeploymentStrategy{},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    nil,
		Paused:                  false,
		ProgressDeadlineSeconds: nil,
	}
	u.deployment = deployment
}

func (u *PodDataService) getContainerPort(podInfo *pod.PodInfo) (containerPort []v13.ContainerPort) {
	for _, v := range podInfo.PodPort {
		containerPort = append(containerPort, v13.ContainerPort{
			Name:          "port-" + strconv.FormatInt(int64(v.ContainerPort), 10),
			ContainerPort: v.ContainerPort,
			Protocol:      u.getProtocol(v.Protocol),
		})
	}
	return
}

func (u *PodDataService) getProtocol(protocol string) v13.Protocol {
	switch protocol {
	case "TCP":
		return "TCP"
	case "UDP":
		return "UDP"
	case "SCTP":
		return "SCTP"
	default:
		return "TCP"
	}
}

func (u *PodDataService) getEnv(podInfo *pod.PodInfo) (envVar []v13.EnvVar) {
	for _, v := range podInfo.PodEnv {
		envVar = append(envVar, v13.EnvVar{
			Name:      v.EnvKey,
			Value:     v.EnvValue,
			ValueFrom: nil,
		})
	}
	return
}

func (u *PodDataService) getResources(podInfo *pod.PodInfo) (source v13.ResourceRequirements) {
	// Maximum resources that can be used
	source.Limits = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(podInfo.PodCpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(podInfo.PodMemoryMax), 'f', 6, 64)),
	}
	// Minimum resources required
	//@TODO Implement dynamic settings
	source.Requests = v13.ResourceList{
		"cpu":    resource.MustParse(strconv.FormatFloat(float64(podInfo.PodCpuMax), 'f', 6, 64)),
		"memory": resource.MustParse(strconv.FormatFloat(float64(podInfo.PodMemoryMax), 'f', 6, 64)),
	}
	return
}

func (u *PodDataService) getImagePullPolicy(podInfo *pod.PodInfo) v13.PullPolicy {
	switch podInfo.PodPullPolicy {
	case "Always":
		return "Always"
	case "Never":
		return "Never"
	case "IfNotPresent":
		return "IfNotPresent"
	default:
		return "Always"
	}
}

// Add Pod
func (u *PodDataService) AddPod(pod2 *model.Pod) (int64, error) {
	return u.PodRepository.CreatePod(pod2)
}

// Delete
func (u *PodDataService) DeletePod(podID int64) error {
	return u.PodRepository.DeletePodByID(podID)
}

// Update
func (u *PodDataService) UpdatePod(pod2 *model.Pod) error {
	return u.PodRepository.UpdatePod(pod2)
}

// Find by ID
func (u *PodDataService) FindPodByID(podID int64) (*model.Pod, error) {
	return u.PodRepository.FindPodByID(podID)
}

// Find all
func (u *PodDataService) FindAllPod() ([]model.Pod, error) {
	return u.PodRepository.FindAll()
}
