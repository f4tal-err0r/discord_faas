package runner

import (
	"net"

	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"

	pb "github.com/f4tal-err0r/discord_faas/proto"
)

type Service struct {
	grpcserv *grpc.Server
	cs       *kubernetes.Clientset
	Proc     *ProcessorService
	pl       Platform
}

type Platform interface {
	NewFuncRunner(img string, funcpath string)
	NewBuilder(img string, funcpath string)
}

func NewService(cs *kubernetes.Clientset) (*Service, error) {
	var svc Service

	svc.Proc = NewProcessorService()
	pb.RegisterProcessorServiceServer(grpc.NewServer(), svc.Proc)

	svc.cs = cs
	return &svc, nil
}

func (s *Service) Build(img string, hash string) error {

}

func (s *Service) Run(lis net.Listener) error {
	return s.grpcserv.Serve(lis)
}
