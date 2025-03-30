package runner

import (
	"net"

	"google.golang.org/grpc"
	spec "k8s.io/api/core/v1"

	pb "github.com/f4tal-err0r/discord_faas/proto"
)

type Service struct {
	grpcserv *grpc.Server
	Proc     *ProcessorService
}

type RunnerOpts struct {
	Id      string
	Image   string
	Cmd     []string
	EnvVars []spec.EnvVar
}

type Platform interface {
	Run() error
	TailLogs() (chan []byte, error)
}

func NewService() (*Service, error) {
	var svc Service

	svc.Proc = NewProcessorService()
	svc.grpcserv = grpc.NewServer()
	pb.RegisterProcessorServiceServer(svc.grpcserv, svc.Proc)

	return &svc, nil
}

func (s *Service) Run(lis net.Listener) error {
	return s.grpcserv.Serve(lis)
}
