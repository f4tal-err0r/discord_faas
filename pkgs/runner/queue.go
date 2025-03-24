package runner

import (
	"context"
	"sync"

	pb "github.com/f4tal-err0r/discord_faas/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ProcessorService struct {
	queue map[string]chan *pb.DiscordContent
	notif *NotifierRouter
	mu    sync.RWMutex
	pb.UnimplementedProcessorServiceServer
}

type NotifierRouter struct {
	queue sync.Map
}

func (n *NotifierRouter) CreateNotifier(funcid string) {
	n.queue.Store(funcid, make(chan *pb.DiscordResp))
}
func NewProcessorService() *ProcessorService {
	return &ProcessorService{
		queue: make(map[string]chan *pb.DiscordContent),
		notif: &NotifierRouter{
			queue: sync.Map{},
		},
	}
}

func (s *ProcessorService) AddContent(c *pb.DiscordContent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.queue[c.Funcmeta.Id]; !exists {
		s.queue[c.Funcmeta.Id] = make(chan *pb.DiscordContent, 10)
	}

	s.queue[c.Funcmeta.Id] <- c
	return nil

}

func (s *ProcessorService) GetWorkerResp(funcid string, msgid string) chan *pb.DiscordResp {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if res, ok := s.notif.queue.Load(msgid); !ok {
		return nil
	} else {
		return res.(chan *pb.DiscordResp)
	}
}

func (s *ProcessorService) SendResp(ctx context.Context, resp *pb.DiscordResp) (*emptypb.Empty, error) {
	s.notif.queue.Store(resp.Funcmeta.Id, resp)
	return &emptypb.Empty{}, nil
}

func (s *ProcessorService) SubContent(funcid *pb.Funcmeta, stream pb.ProcessorService_RecvContentServer) error {
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
