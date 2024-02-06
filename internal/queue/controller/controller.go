package controller

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/service"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/queue/provider"
)

type Controller struct {
	cache *Cache[*pb.QueueItem]
}

func New(ctx context.Context, d service.GrpcDialer, o Option) (*Controller, error) {
	c, err := NewCache[*pb.QueueItem](ctx, d, o)
	if err != nil {
		return nil, err
	}
	return &Controller{
		cache: c,
	}, nil
}

func (s *Controller) Queue(ctx context.Context, in *pb.QueueRequest) (*pb.Void, error) {
	myqueue, err := s.cache.GetQueue(ctx, in.Tenant, in.Queue)
	if err != nil {
		return nil, err
	}
	for _, i := range in.Items {
		if err := myqueue.Add(i); err != nil {
			return nil, err
		}
	}
	ret := pb.Void{}
	return &ret, nil
}

func (s *Controller) Close() error {
	return s.cache.Close()
}

func (s *Controller) Dequeue(ctx context.Context, in *pb.DequeueRequest) (*pb.DequeueReply, error) {
	myqueue, err := s.cache.GetQueue(ctx, in.Tenant, in.Queue)
	if err != nil {
		return nil, err
	}
	i, err := myqueue.Remove()
	if err != nil && !errors.Is(err, provider.ErrQueueEmpty) {
		return nil, err
	}
	var iqs []*pb.QueueItem
	if i != nil {
		iqs = append(iqs, i)
	}
	return &pb.DequeueReply{
		Items: iqs,
	}, nil
}
