package handler

import (
	"context"
	"git.imooc.com/coding-535/common"
	"git.imooc.com/coding-535/volume/domain/model"
	"git.imooc.com/coding-535/volume/domain/service"
	log "github.com/asim/go-micro/v3/logger"
	volume "git.imooc.com/coding-535/volume/proto/volume"
	"strconv"
)

type VolumeHandler struct{
     // Note: This type implements the IVolumeDataService interface
     VolumeDataService service.IVolumeDataService
}

// Call is a single request handler called via client.Call or the generated client code
func (e *VolumeHandler) AddVolume(ctx context.Context,info *volume.VolumeInfo , rsp *volume.Response) error {
	log.Info("Received *volume.AddVolume request")
	volume := &model.Volume{}
	if err := common.SwapTo(info,volume);err!=nil{
		common.Error(err)
		rsp.Msg= err.Error()
		return err
	}
	// Create volume
	if err:= e.VolumeDataService.CreateVolumeToK8s(info);err!=nil{
		common.Error(err)
		rsp.Msg= err.Error()
		return err
	} else {
		// Write to database
		volumeID,err := e.VolumeDataService.AddVolume(volume)
		if err != nil {
			common.Error(err)
			rsp.Msg= err.Error()
			return err
		}
		rsp.Msg = "Volume added successfully, ID: "+ strconv.FormatInt(volumeID,10)
	}
	return nil
}

// Delete
func (e *VolumeHandler) DeleteVolume(ctx context.Context, req *volume.VolumeId, rsp *volume.Response) error {
	log.Info("Received *volume.DeleteVolume request")
	volumModel,err := e.VolumeDataService.FindVolumeByID(req.Id)
	if err != nil {
		common.Error(err)
		return err
	}
	// Delete from k8s and database
	if err := e.VolumeDataService.DeleteVolumeFromK8s(volumModel);err !=nil{
		common.Error(err)
		return err
	}
	return nil
}

func (e *VolumeHandler) UpdateVolume(ctx context.Context, req *volume.VolumeInfo, rsp *volume.Response) error {
	log.Info("Received *volume.UpdateVolume request")
	return nil
}

// Find volume by ID
func (e *VolumeHandler) FindVolumeByID(ctx context.Context, req *volume.VolumeId, rsp *volume.VolumeInfo) error {
	log.Info("Received *volume.FindVolumeByID request")
	volumeModel,err:=e.VolumeDataService.FindVolumeByID(req.Id)
	if err != nil {
		common.Error(err)
		return err
	}
	// Data conversion
	if err := common.SwapTo(volumeModel,rsp);err!=nil{
		common.Error(err)
		return err
	}
	return nil
}

func (e *VolumeHandler) FindAllVolume(ctx context.Context, req *volume.FindAll, rsp *volume.AllVolume) error {
	log.Info("Received *volume.FindAllVolume request")
	allVolume,err := e.VolumeDataService.FindAllVolume()
	if err != nil {
		common.Error(err)
		return err
	}
	// Format response
	for _,v :=range  allVolume{
		// Create instance
		volumeInfo := &volume.VolumeInfo{}
		// Data conversion
		if err:= common.SwapTo(v,volumeInfo);err !=nil{
			common.Error(err)
			return err
		}
		// Merge data
		rsp.VolumeInfo = append(rsp.VolumeInfo,volumeInfo)
	}
	return nil
}


