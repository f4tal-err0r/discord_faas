package runner

import (
	"net"

	"google.golang.org/grpc"

	pb "github.com/f4tal-err0r/discord_faas/proto"
)

type Service struct {
	grpcserv *grpc.Server
	Proc     *ProcessorService
}

type RunnerOpts struct {
	Id    string
	Image string
	Name  string
	Cmd   []string
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

func (s *Service) CreateRunner(pl Platform, content *pb.DiscordContent) {
	s.Proc.AddContent()
}
