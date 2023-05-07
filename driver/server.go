package driver

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"github.com/terloo/nfs-csi/csi"
	"os"
	"strings"
)

type nonBlockGPRCServer struct {
	Server *grpc.Server
}

func NewNonBlockGRPCServer() *nonBlockGPRCServer {
	return &nonBlockGPRCServer{Server: grpc.NewServer()}
}

func (s *nonBlockGPRCServer) RunControllerServer(endpoint string) {
	proto, addr, err := ParseEndpoint(endpoint)
	if err != nil {
		panic(err)
	}

	csi.RegisterIdentityServer(s.Server, NewIdentityServer())
	csi.RegisterControllerServer(s.Server, NewControllerServer())

	if strings.HasPrefix(endpoint, "unix://") {
		os.Remove(endpoint)
	}

	listen, err := net.Listen(proto, addr)
	if err != nil {
		panic(err)
	}

	reflection.Register(s.Server)
	log.Println("Server running in " + endpoint + " ....")
	err = s.Server.Serve(listen)
	if err != nil {
		panic(err)
	}

}

func (s *nonBlockGPRCServer) RunNodeServer(endpoint string) {
	proto, addr, err := ParseEndpoint(endpoint)
	if err != nil {
		panic(err)
	}

	csi.RegisterIdentityServer(s.Server, NewIdentityServer())
	csi.RegisterNodeServer(s.Server, NewNodeServer())

	if strings.HasPrefix(endpoint, "unix://") {
		err := os.Remove(addr)
		if err != nil {
			log.Println("删除 socket 失败", err)
		}
	}

	listen, err := net.Listen(proto, addr)
	if err != nil {
		panic(err)
	}

	reflection.Register(s.Server)
	log.Println("Server running in " + endpoint + " ....")
	err = s.Server.Serve(listen)
	if err != nil {
		panic(err)
	}

}