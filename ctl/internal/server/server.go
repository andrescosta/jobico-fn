package server

import (
	"context"
	"strconv"

	"github.com/andrescosta/goico/pkg/database"
	"github.com/andrescosta/goico/pkg/reflectico"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/ctl/internal/dao"
	"google.golang.org/protobuf/proto"
)

type Queue struct {
	pb.UnimplementedControlServer
	daos map[string]*dao.ProtoMessageDAO
	db   *database.Database
}

const DB_PATH = ".\\db.db"

func NewQueue(ctx context.Context) (*Queue, error) {
	q, err := database.Open(ctx, DB_PATH)
	if err != nil {
		return nil, err
	}
	return &Queue{
		daos: make(map[string]*dao.ProtoMessageDAO),
		db:   q,
	}, nil
}

func (s *Queue) GetQueues(ctx context.Context, in *pb.GetQueuesDefRequest) (*pb.GetQueuesDefReply, error) {
	r := &pb.GetQueuesDefReply{}
	mydao, err := s.getDao(ctx, in.MerchantId.Id, "queue", &pb.QueueDef{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.All(ctx)
	if err != nil {
		return nil, err
	}
	queues := reflectico.Convert[*proto.Message, *pb.QueueDef](ms)
	r.Queues = queues

	return r, nil
}

func (s *Queue) AddQueues(ctx context.Context, in *pb.AddQueueDefRequest) (*pb.AddQueueDefReply, error) {
	mydao, err := s.getDao(ctx, in.Queues.MerchantId.Id, "queue", &pb.QueueDef{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Queues
	ms, err := mydao.Add(ctx, &m)
	if err != nil {
		return nil, err
	}
	in.Queues.ID = strconv.FormatUint(ms, 10)
	return &pb.AddQueueDefReply{
		Queues: in.Queues,
	}, nil

}

func (s *Queue) GetEvents(ctx context.Context, in *pb.GetEventRequest) (*pb.GetEventReply, error) {
	return nil, nil
}
func (s *Queue) AddEvents(ctx context.Context, in *pb.AddEventRequest) (*pb.AddEventReply, error) {
	return nil, nil
}
func (s *Queue) GetListeners(ctx context.Context, in *pb.GetListenersRequest) (*pb.GetListenersReply, error) {
	return nil, nil
}
func (s *Queue) AddListener(ctx context.Context, in *pb.AddListenersRequest) (*pb.AddListenersReply, error) {
	return nil, nil
}
func (s *Queue) AddExecutor(ctx context.Context, in *pb.AddExecutorRequest) (*pb.AddExecutorReply, error) {
	return nil, nil
}
func (s *Queue) GetExecutors(ctx context.Context, in *pb.GetExecutorsRequest) (*pb.GetExecutorsReply, error) {
	return nil, nil
}

func (s *Queue) getDao(ctx context.Context, merchant string, entity string, message proto.Message) (*dao.ProtoMessageDAO, error) {
	mydao, ok := s.daos[merchant]
	if !ok {
		var err error
		mydao, err = dao.NewProtoMessageDAO(ctx, s.db, merchant+"/"+entity, message)
		if err != nil {
			return nil, err
		}
		s.daos[merchant] = mydao
	}
	return mydao, nil
}
