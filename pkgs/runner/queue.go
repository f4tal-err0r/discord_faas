package runner

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"

	pb "github.com/f4tal-err0r/discord_faas/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ProcessorService struct {
	notif notifyQueue
	queue map[string]chan *pb.DiscordContent
	mu    sync.RWMutex
	pb.UnimplementedProcessorServiceServer
}

type notifyQueue struct {
	queue map[string]chan *pb.DiscordResp
	mu    sync.RWMutex
}

func NewProcessorService() *ProcessorService {
	return &ProcessorService{
		queue: make(map[string]chan *pb.DiscordContent),
		notif: notifyQueue{
			queue: make(map[string]chan *pb.DiscordResp),
		},
	}
}

func (s *ProcessorService) AddContent(c *pb.DiscordContent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.queue[c.Msgid]; !exists {
		s.queue[c.Msgid] = make(chan *pb.DiscordContent, 10)
	}

	s.queue[c.Msgid] <- c
	return nil

}

func (s *notifyQueue) GetAddWorkerQueue(id string) chan *pb.DiscordResp {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.queue[id]; !ok {
		s.queue[id] = make(chan *pb.DiscordResp, 10)
		return s.queue[id]
	} else {
		return c
	}
}

func (s *ProcessorService) SendResp(ctx context.Context, resp *pb.DiscordResp) (*emptypb.Empty, error) {
	s.notif.queue[resp.Msgid] <- resp
	return &emptypb.Empty{}, nil
}

func (s *ProcessorService) SubContent(msgID *pb.Workloadid, stream pb.ProcessorService_RecvContentServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if c, ok := s.queue[msgID.Id]; !ok {
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

func generateHex() string {
	b := make([]byte, 5) // 5 bytes = 10 hex characters
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
