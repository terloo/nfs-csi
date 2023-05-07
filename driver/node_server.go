package driver

import (
	"context"
	"errors"
	"log"
	"github.com/terloo/nfs-csi/csi"
	"os"
)

var _ csi.NodeServer = &NodeServer{}

type NodeServer struct {
	csi.UnimplementedNodeServer
}

func NewNodeServer() *NodeServer {
	return &NodeServer{}
}

func (n NodeServer) NodeGetCapabilities(ctx context.Context, request *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{Capabilities: []*csi.NodeServiceCapability{
		{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
				},
			},
		},
	}}, nil
}

func (n NodeServer) NodeGetInfo(ctx context.Context, request *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId:            "csi-master",
		MaxVolumesPerNode: 100,
	}, nil
}

// 对容器卷进行格式化，并将其挂载到临时目录(暂存)
func (n NodeServer) NodeStageVolume(ctx context.Context, request *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	// nfs不需要暂存
	log.Printf("NodeStageVolume %v\n", request)
	log.Printf("暂存一个底层存储%s到目录%s", request.VolumeId, request.StagingTargetPath)
	return &csi.NodeStageVolumeResponse{}, nil
}

func (n NodeServer) NodeUnstageVolume(ctx context.Context, request *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	log.Printf("NodeUnstageVolume %v\n", request)
	return &csi.NodeUnstageVolumeResponse{}, nil
}

// 将容器卷挂载到容器对应目录
func (n NodeServer) NodePublishVolume(ctx context.Context, request *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	log.Printf("NodePublishVolume %v\n", request)
	log.Printf("推送容器卷%s到目录%s", request.VolumeId, request.TargetPath)
	targetPath := request.TargetPath

	err := os.MkdirAll(targetPath, os.ModePerm)
	if err != nil {
		return nil, errors.New("创建容器目录失败" + err.Error())
	}

	if _, err := os.Lstat(targetPath); !os.IsNotExist(err) {
		log.Println("已经挂载目录", targetPath)
		return &csi.NodePublishVolumeResponse{}, nil
	}

	// 挂载nfs
	nfsVolumeDir, ok := request.PublishContext["nfsVolumeDir"]
	if !ok {
		return nil, errors.New("无法获取nfs服务器PV路径")
	}
	log.Println("假装挂载到了", nfsVolumeDir)
	//command := exec.Command("mount", "-t", "nfs", "10.0.0.9:"+nfsVolumeDir, targetPath)
	//output, err := command.Output()
	//if err != nil {
	//	return nil, errors.New("执行挂载命令失败" + string(output))
	//}
	log.Println("推送成功")
	log.Println()
	return &csi.NodePublishVolumeResponse{}, nil
}

func (n NodeServer) NodeUnpublishVolume(ctx context.Context, request *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	log.Printf("NodeUnpublishVolume %v\n", request)
	log.Printf("取消推送容器卷%s到目录%s", request.VolumeId, request.TargetPath)
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (n NodeServer) NodeGetVolumeStats(ctx context.Context, request *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	// TODO implement me
	panic("implement NodeGetVolumeStats")
}

func (n NodeServer) NodeExpandVolume(ctx context.Context, request *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	// TODO implement me
	panic("implement NodeExpandVolume")
}
