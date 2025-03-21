package runner

import (
	"net"

	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	pb "github.com/f4tal-err0r/discord_faas/proto"
)

type Service struct {
	grpcserv *grpc.Server
	cs       *kubernetes.Clientset
	Proc     *ProcessorService
	platform
}

type Platform interface {
	NewFuncRunner(img string, funcpath string)
	NewBuilder(img string, funcpath string)
}

func NewService() (*Service, error) {
	var svc Service

	svc.Proc = NewProcessorService()
	pb.RegisterProcessorServiceServer(grpc.NewServer(), svc.Proc)

	//access internal k8s client
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	svc.cs = cs
	return &svc, nil
}

func (s *Service) Build(pl Platform, img string, hash string) error {

}

func (s *Service) Run(lis net.Listener) error {
	return s.grpcserv.Serve(lis)
}
