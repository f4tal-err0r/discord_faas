package server

import (
	"context"
	"sync"

	pb "github.com/f4tal-err0r/discord_faas/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type FuncService struct {
	queue []*pb.DiscordContent
	mu    sync.RWMutex
	pb.UnimplementedProcessorServiceServer
}

func NewProcessorService() *FuncService {
	return &FuncService{
		queue: []*pb.DiscordContent{},
	}
}

func (s *FuncService) AddContent(c *pb.DiscordContent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.queue = append(s.queue, c)
	return nil

}

func (s *FuncService) RecvContent(ctx context.Context, stream pb.ProcessorService_RecvContentServer) error {

}

func (s *FuncService) SendResp(ctx context.Context, resp *pb.DiscordResp) (*emptypb.Empty, error) {
	s.notif.queue.Store(resp.Funcmeta.Id, resp)
	return &emptypb.Empty{}, nil
}

func (s *FuncService) SubContent(funcid *pb.Funcmeta, stream pb.ProcessorService_RecvContentServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if c, ok := s.queue[funcid.Id]; !ok {
		return nil
	} else {
		for content := range c {
			if err := stream.Send(content); err != nil {
				return err
			}
		}
		return nil
	}
}
