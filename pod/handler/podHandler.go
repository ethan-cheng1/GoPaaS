package handler

import (
	"context"
	"git.imooc.com/coding-535/common"
	"git.imooc.com/coding-535/pod/domain/model"
	"git.imooc.com/coding-535/pod/domain/service"
	"git.imooc.com/coding-535/pod/proto/pod"
	"strconv"
)

type PodHandler struct {
	PodDataService service.IPodDataService
}

// Add and create POD
func (e *PodHandler) AddPod(ctx context.Context, info *pod.PodInfo, rsp *pod.Response) error {
	common.Info("Adding pod")
	podModel := &model.Pod{}
	err := common.SwapTo(info,podModel)
	if err != nil {
		common.Error(err)
		rsp.Msg = err.Error()
		return err
	}

	if err:=e.PodDataService.CreateToK8s(info);err != nil {
		common.Error(err)
		rsp.Msg = err.Error()
		return err
	} else {
		// Write data to database
		podID,err := e.PodDataService.AddPod(podModel)
		if err != nil {
			common.Error(err)
			rsp.Msg = err.Error()
			return err
		}
		common.Info("Pod added successfully, database ID: " + strconv.FormatInt(podID, 10))
		rsp.Msg = "Pod added successfully, database ID: " + strconv.FormatInt(podID, 10)
	}
	return nil
}

// Delete pod from k8s and database
func (e *PodHandler) DeletePod(ctx context.Context, req *pod.PodId, rsp *pod.Response) error {
	// Find data
	podModel ,err := e.PodDataService.FindPodByID(req.Id)
	if err != nil {
		common.Error(err)
		return err
	}
	if err := e.PodDataService.DeleteFromK8s(podModel);err!=nil{
		common.Error(err)
		return err
	}
	return nil
}

// Update given pod
func (e *PodHandler) UpdatePod(ctx context.Context, req *pod.PodInfo, rsp *pod.Response) error {
	// Update pod information in k8s
	err := e.PodDataService.UpdateToK8s(req)
	if err != nil {
		common.Error(err)
		return err
	}
	// Query pod in database
	podModel,err:=e.PodDataService.FindPodByID(req.Id)
	if err != nil {
		common.Error(err)
		return err
	}
	err = common.SwapTo(req,podModel)
	if err != nil {
		common.Error(err)
		return err
	}
	e.PodDataService.UpdatePod(podModel)
	return nil

}

// Query single pod information
func (e *PodHandler) FindPodByID(ctx context.Context,req *pod.PodId,rsp *pod.PodInfo) error  {
	// Query pod data
	podModel ,err := e.PodDataService.FindPodByID(req.Id)
	if err != nil {
		common.Error(err)
		return err
	}
	err = common.SwapTo(podModel,rsp)
	if err != nil {
		common.Error(err)
		return err
	}
	return nil

}

// Query all pods
func (e *PodHandler) FindAllPod(ctx context.Context, req *pod.FindAll, rsp *pod.AllPod) error {
	// Query all pods
	allPod ,err := e.PodDataService.FindAllPod()
	if err != nil {
		common.Error(err)
		return err
	}
	// Format the response
	for _,v:=range allPod{
		podInfo := &pod.PodInfo{}
		err := common.SwapTo(v,podInfo)
		if err != nil {
			common.Error(err)
			return err
		}
		rsp.PodInfo = append(rsp.PodInfo,podInfo)
	}
	return nil
}



