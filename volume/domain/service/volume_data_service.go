package service

import (
	"context"
	"errors"
	"git.imooc.com/coding-535/common"
	"git.imooc.com/coding-535/volume/domain/model"
	"git.imooc.com/coding-535/volume/domain/repository"
	"git.imooc.com/coding-535/volume/proto/volume"
	"k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

// Interface type definition
type IVolumeDataService interface {
	AddVolume(*model.Volume) (int64 , error)
	DeleteVolume(int64) error
	UpdateVolume(*model.Volume) error
	FindVolumeByID(int64) (*model.Volume, error)
	FindAllVolume() ([]model.Volume, error)

	CreateVolumeToK8s(*volume.VolumeInfo) error
	DeleteVolumeFromK8s(*model.Volume) error
}

// Create
// Note: Returns IVolumeDataService interface type
func NewVolumeDataService(volumeRepository repository.IVolumeRepository,clientSet *kubernetes.Clientset) IVolumeDataService{
	return &VolumeDataService{ VolumeRepository:volumeRepository, K8sClientSet: clientSet,deployment:&v1.Deployment{}}
}

type VolumeDataService struct {
    // Note: This is IVolumeRepository type
	VolumeRepository repository.IVolumeRepository
	K8sClientSet  *kubernetes.Clientset
	deployment  *v1.Deployment
}

// Delete PVC from k8s
func (u *VolumeDataService) DeleteVolumeFromK8s(volume *model.Volume) (err error) {
	// Delete from K8s
	if err = u.K8sClientSet.CoreV1().PersistentVolumeClaims(volume.VolumeNamespace).Delete(context.TODO(),volume.VolumeName, v13.DeleteOptions{});err!=nil {
		common.Error(err)
		return err
	} else {
		// Delete from database
		if err := u.DeleteVolume(volume.ID);err !=nil {
			common.Error(err)
			return err
		}
		common.Info("Storage ID " + strconv.FormatInt(volume.ID,10) + " deleted successfully!")
	}
	return
}

// Create storage in k8s
func (u *VolumeDataService) CreateVolumeToK8s(info *volume.VolumeInfo)(err error) {
	volume := u.setVolume(info)
	if _,err = u.K8sClientSet.CoreV1().PersistentVolumeClaims(info.VolumeNamespace).Get(context.TODO(),info.VolumeName,v13.GetOptions{});err !=nil {
		// If storage does not exist, create it
		if _,err = u.K8sClientSet.CoreV1().PersistentVolumeClaims(info.VolumeNamespace).Create(context.TODO(),volume,v13.CreateOptions{});err !=nil{
			common.Error(err)
			return err
		}
		common.Info("Storage created successfully")
		return nil
	}else {
		common.Error("Storage space " + info.VolumeName + " already exists")
		return errors.New("Storage space " + info.VolumeName + " already exists")
	}
}

// Set PVC details
func (u *VolumeDataService) setVolume(info *volume.VolumeInfo) *v12.PersistentVolumeClaim {
	pvc := &v12.PersistentVolumeClaim{}
	// Set interface type
	pvc.TypeMeta = v13.TypeMeta{
		Kind:       "PersistentVolumeClaim",
		APIVersion: "v1",
	}
	// Set storage basic information
	pvc.ObjectMeta = v13.ObjectMeta{
		Name:                       info.VolumeName,
		Namespace:                  info.VolumeNamespace,
		Annotations: map[string]string{
			"pv.kubernetes.io/bound-by-controller":"yes",
			"volume.beta.kubernetes.io/storage-provisioner":"rbd.csi.ceph.com",
			"Cap":"GoPaaS Platform",
		},
	}
	// Set storage dynamic information
	pvc.Spec = v12.PersistentVolumeClaimSpec{
		AccessModes:      u.getAccessModes(info),
		Resources:        u.getResource(info),
		StorageClassName: &info.VolumeStorageClassName,
		VolumeMode:       u.getVolumeMode(info),
	}
	return pvc

}

// Get storage type
func (u *VolumeDataService) getVolumeMode(info *volume.VolumeInfo) *v12.PersistentVolumeMode {
	var pvm v12.PersistentVolumeMode
	switch info.VolumePersistentVolumeMode {
	case "Block":
		pvm = v12.PersistentVolumeBlock
	case "Filesystem":
		pvm = v12.PersistentVolumeFilesystem
	default:
		pvm = v12.PersistentVolumeFilesystem
	}
	return &pvm
}

// Get resource configuration
func (u *VolumeDataService) getResource (info *volume.VolumeInfo)(source v12.ResourceRequirements)  {
	source.Requests = v12.ResourceList{
		"storage": resource.MustParse(strconv.FormatFloat(float64(info.VolumeRequest),'f',6,64)+"Gi"),
	}
	return
}

// Get access modes
func (u *VolumeDataService) getAccessModes(info *volume.VolumeInfo)(pvam []v12.PersistentVolumeAccessMode)  {
	var pm v12.PersistentVolumeAccessMode
	switch info.VolumeAccessMode {
	case "ReadWriteOnce":
		pm = v12.ReadWriteOnce
	case "ReadOnlyMany":
		pm = v12.ReadOnlyMany
	case "ReadWriteMany":
		pm = v12.ReadWriteMany
	case "ReadWriteOncePod":
		pm = v12.ReadWriteOncePod
	default:
		pm = v12.ReadWriteOnce
	}
	pvam = append(pvam,pm)
	return pvam
	
}

// Insert
func (u *VolumeDataService) AddVolume(volume *model.Volume) (int64 ,error) {
	 return u.VolumeRepository.CreateVolume(volume)
}

// Delete
func (u *VolumeDataService) DeleteVolume(volumeID int64) error {
	return u.VolumeRepository.DeleteVolumeByID(volumeID)
}

// Update
func (u *VolumeDataService) UpdateVolume(volume *model.Volume) error {
	return u.VolumeRepository.UpdateVolume(volume)
}

// Find
func (u *VolumeDataService) FindVolumeByID(volumeID int64) (*model.Volume, error) {
	return u.VolumeRepository.FindVolumeByID(volumeID)
}

// Find all
func (u *VolumeDataService) FindAllVolume() ([]model.Volume, error) {
	return u.VolumeRepository.FindAll()
}

