package server

import (
	"context"

	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/ctl/internal/dao"
)

type Queue struct {
	pb.UnimplementedControlServer
	daos map[string]*dao.QueueDAO
}

func NewQueue() *Queue {
	return &Queue{
		daos: make(map[string]*dao.QueueDAO),
	}
}

func (s *Queue) GetQueues(ctx context.Context, in *pb.GetQueuesDefRequest) (*pb.GetQueuesDefReply, error) {
	r := &pb.GetQueuesDefReply{}
	mydao, err := s.getDao(ctx, in.MerchantId.Id)
	if err != nil {
		return nil, err
	}
	qs, err := mydao.GetQueueDefs(ctx)
	if err != nil {
		return nil, err
	}
	r.Queues = qs
	return r, nil
}

func (s *Queue) AddQueues(ctx context.Context, in *pb.AddQueueDefRequest) (*pb.AddQueueDefReply, error) {
	mydao, err := s.getDao(ctx, in.Queues.MerchantId.Id)
	if err != nil {
		return nil, err
	}
	r := &pb.AddQueueDefReply{}
	_, err = mydao.AddQueueDef(ctx, in.Queues)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Queue) getDao(ctx context.Context, merchant string) (*dao.QueueDAO, error) {
	mydao, ok := s.daos[merchant]
	if !ok {
		var err error
		mydao, err = dao.NewQueueDAO(ctx, ".\\d.db")
		if err != nil {
			return nil, err
		}
		s.daos[merchant] = mydao
	}
	return mydao, nil
}
