package driver

import (
	"context"
	"github.com/terloo/nfs-csi/csi"
	"log"
	"os"
)

var _ csi.ControllerServer = &ControllerServer{}

type ControllerServer struct {
	csi.UnimplementedControllerServer
}

func NewControllerServer() *ControllerServer {
	return &ControllerServer{}
}

func (c ControllerServer) ControllerGetCapabilities(ctx context.Context, request *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return &csi.ControllerGetCapabilitiesResponse{Capabilities: []*csi.ControllerServiceCapability{
		{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME}},
		},
		{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME}},
		},
	}}, nil
}

// 创建底层存储
func (c ControllerServer) CreateVolume(ctx context.Context, request *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	log.Printf("CreateVolume %v\n", request)
	log.Printf("创建一个容器卷，名字为%s，类型为%v", request.Name, request.VolumeCapabilities)

	// 在nfs服务器上创建一个目录
	nfsVolumeDir := "/nfs/data/" + request.Name
	err := os.Mkdir(nfsVolumeDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{VolumeId: request.Name},
	}, nil
}

// 删除底层存储
func (c ControllerServer) DeleteVolume(ctx context.Context, request *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	log.Printf("DeleteVolume %v\n", request)

	// 删除nfs服务器上的目录
	nfsPVCDir := "/nfs/data/" + request.VolumeId
	err := os.RemoveAll(nfsPVCDir)
	if err != nil {
		return nil, err
	}

	return &csi.DeleteVolumeResponse{}, nil
}

// 将底层存储推送到目标节点上
func (c ControllerServer) ControllerPublishVolume(ctx context.Context, request *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	log.Printf("ControllerPublishVolume %v\n", request)
	log.Printf("推送一个容器卷%s到节点%v", request.VolumeId, request.NodeId)
	nfsVolumeDir := "/nfs/data/" + request.VolumeId
	return &csi.ControllerPublishVolumeResponse{PublishContext: map[string]string{
		"nfsVolumeDir": nfsVolumeDir,
	}}, nil
}

// 将底层存储从目标节点移除
func (c ControllerServer) ControllerUnpublishVolume(ctx context.Context, request *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	log.Printf("ControllerUnpublishVolume %v\n", request)
	log.Printf("从节点%v取消推送一个容器卷%s", request.NodeId, request.VolumeId)
	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

func (c ControllerServer) ValidateVolumeCapabilities(ctx context.Context, request *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	// TODO implement me
	panic("implement ValidateVolumeCapabilities")
}

func (c ControllerServer) ListVolumes(ctx context.Context, request *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	// TODO implement me
	panic("implement ListVolumes")
}

func (c ControllerServer) GetCapacity(ctx context.Context, request *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	// TODO implement me
	panic("implement GetCapacity")
}

func (c ControllerServer) CreateSnapshot(ctx context.Context, request *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	// TODO implement me
	panic("implement CreateSnapshot")
}

func (c ControllerServer) DeleteSnapshot(ctx context.Context, request *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	// TODO implement me
	panic("implement DeleteSnapshot")
}

func (c ControllerServer) ListSnapshots(ctx context.Context, request *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	// TODO implement me
	panic("implement ListSnapshots")
}

func (c ControllerServer) ControllerExpandVolume(ctx context.Context, request *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	// TODO implement me
	panic("implement ControllerExpandVolume")
}

func (c ControllerServer) ControllerGetVolume(ctx context.Context, request *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	// TODO implement me
	panic("implement ControllerGetVolume")
}
