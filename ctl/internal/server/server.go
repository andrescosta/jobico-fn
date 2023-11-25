package server

import (
	"context"
	"fmt"

	"github.com/andrescosta/goico/pkg/convertico"
	"github.com/andrescosta/goico/pkg/database"
	pb "github.com/andrescosta/workflew/api/types"
	"github.com/andrescosta/workflew/ctl/internal/dao"
	"google.golang.org/protobuf/proto"
)

const (
	DB_PATH         = ".\\db.db"
	TBL_EVENT       = "event"
	TBL_QUEUE       = "queue"
	TBL_PACKAGE     = "package"
	TBL_LISTENER    = "listener"
	TBL_MERCHANT    = "merchant"
	TBL_ENVIRONMENT = "environment"
	TBL_EXECUTOR    = "executor"
	GEN_MERCHANT    = "[Generic]"
)

type ControlServer struct {
	pb.UnimplementedControlServer
	daos map[string]*dao.DAO[proto.Message]
	db   *database.Database
}

func NewCotrolServer(ctx context.Context) (*ControlServer, error) {
	db, err := database.Open(ctx, DB_PATH)
	if err != nil {
		return nil, err
	}
	return &ControlServer{
		daos: make(map[string]*dao.DAO[proto.Message]),
		db:   db,
	}, nil
}

func (s *ControlServer) Close(ctx context.Context) {
	s.db.Close(ctx)
}

func (s *ControlServer) GetPackages(ctx context.Context, in *pb.GetJobPackagesRequest) (*pb.GetJobPackagesReply, error) {
	mydao, err := s.getDao(ctx, in.MerchantId, TBL_PACKAGE, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.All(ctx)
	if err != nil {
		return nil, err
	}
	packages := convertico.ConvertSlices[proto.Message, *pb.JobPackage](ms)

	return &pb.GetJobPackagesReply{Packages: packages}, nil
}

func (s *ControlServer) GetAllPackages(ctx context.Context, in *pb.GetAllJobPackagesRequest) (*pb.GetAllJobPackagesReply, error) {
	ms, err := s.getMechants(ctx)
	if err != nil {
		return nil, err
	}
	packages := make([]*pb.JobPackage, 0)
	for _, me := range ms {
		mydao, err := s.getDao(ctx, me.MerchantId, TBL_PACKAGE, &pb.JobPackage{})
		if err != nil {
			return nil, err
		}

		ms, err := mydao.All(ctx)
		if err != nil {
			return nil, err
		}
		ps := convertico.ConvertSlices[proto.Message, *pb.JobPackage](ms)
		packages = append(packages, ps...)
	}
	return &pb.GetAllJobPackagesReply{Packages: packages}, nil
}

func (s *ControlServer) AddPackage(ctx context.Context, in *pb.AddJobPackageRequest) (*pb.AddJobPackageReply, error) {
	mydao, err := s.getDao(ctx, in.Package.MerchantId, TBL_PACKAGE, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Package
	ms, err := mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	in.Package.ID = &ms
	return &pb.AddJobPackageReply{Package: in.Package}, nil

}

func (s *ControlServer) GetMechants(ctx context.Context, in *pb.GetMerchantsRequest) (*pb.GetMerchantsReply, error) {
	m, err := s.getMechants(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetMerchantsReply{Merchants: m}, nil
}

func (s *ControlServer) getMechants(ctx context.Context) ([]*pb.Merchant, error) {
	mydao, err := s.getDaoGen(ctx, TBL_MERCHANT, &pb.Merchant{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.All(ctx)
	if err != nil {
		return nil, err
	}
	mechants := convertico.ConvertSlices[proto.Message, *pb.Merchant](ms)

	return mechants, nil
}

func (s *ControlServer) AddMerchant(ctx context.Context, in *pb.AddMerchantRequest) (*pb.AddMerchantReply, error) {
	mydao, err := s.getDaoGen(ctx, TBL_MERCHANT, &pb.Merchant{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Merchant
	ms, err := mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	in.Merchant.ID = &ms
	return &pb.AddMerchantReply{Merchant: in.Merchant}, nil
}

func (s *ControlServer) AddEnviroment(ctx context.Context, in *pb.AddEnviromentRequest) (*pb.AddEnviromentReply, error) {
	mydao, err := s.getDaoGen(ctx, TBL_ENVIRONMENT, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Environment
	ms, err := mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	in.Environment.ID = &ms
	return &pb.AddEnviromentReply{Environment: in.Environment}, nil
}
func (s *ControlServer) UpdateEnviroment(ctx context.Context, in *pb.UpdateEnviromentRequest) (*pb.UpdateEnviromentReply, error) {
	if in.Environment.ID == nil {
		return nil, fmt.Errorf("ID is empty")
	}
	mydao, err := s.getDaoGen(ctx, TBL_ENVIRONMENT, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Environment
	err = mydao.Update(ctx, m)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateEnviromentReply{}, nil
}

func (s *ControlServer) GetEnviroment(ctx context.Context, in *pb.GetEnviromentRequest) (*pb.GetEnviromentReply, error) {
	mydao, err := s.getDaoGen(ctx, TBL_ENVIRONMENT, &pb.Environment{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.All(ctx)
	if err != nil {
		return nil, err
	}
	environment := &pb.Environment{}
	if len(ms) > 0 {
		environment = ms[0].(*pb.Environment)
	}
	return &pb.GetEnviromentReply{Environment: environment}, nil
}

func (s *ControlServer) getDaoGen(ctx context.Context, entity string, message proto.Message) (*dao.DAO[proto.Message], error) {
	return s.getDao(ctx, GEN_MERCHANT, entity, message)
}

func (s *ControlServer) getDao(ctx context.Context, merchant string, entity string, message proto.Message) (*dao.DAO[proto.Message], error) {
	mydao, ok := s.daos[merchant]
	if !ok {
		var err error
		mydao, err = dao.NewDAO(ctx, s.db, merchant+"/"+entity,
			&ProtoMessageMarshaler{
				newMessage: func() proto.Message {
					return message
				},
			})
		if err != nil {
			return nil, err
		}
		s.daos[merchant] = mydao
	}
	return mydao, nil
}

/*
	func (s *ControlServer) GetQueueDefs(ctx context.Context, in *pb.GetQueueDefsRequest) (*pb.GetQueueDefsReply, error) {
		mydao, err := s.getDao(ctx, in.MerchantId.Id, TBL_QUEUE, &pb.QueueDef{})
		if err != nil {
			return nil, err
		}

		ms, err := mydao.All(ctx)
		if err != nil {
			return nil, err
		}
		queues := convertico.ConvertSlices[proto.Message, *pb.QueueDef](ms)

		return &pb.GetQueueDefsReply{Queues: queues}, nil
	}

	func (s *ControlServer) AddQueueDef(ctx context.Context, in *pb.AddQueueDefRequest) (*pb.AddQueueDefReply, error) {
		mydao, err := s.getDao(ctx, in.Queue.MerchantId.Id, TBL_QUEUE, &pb.QueueDef{})
		if err != nil {
			return nil, err
		}
		var m proto.Message = in.Queue
		ms, err := mydao.Add(ctx, m)
		if err != nil {
			return nil, err
		}
		in.Queue.ID = strconv.FormatUint(ms, 10)
		return &pb.AddQueueDefReply{Queue: in.Queue}, nil

}

	func (s *ControlServer) GetEventDefs(ctx context.Context, in *pb.GetEventDefsRequest) (*pb.GetEventDefsReply, error) {
		mydao, err := s.getDao(ctx, in.MerchantId.Id, TBL_EVENT, &pb.EventDef{})
		if err != nil {
			return nil, err
		}

		ms, err := mydao.All(ctx)
		if err != nil {
			return nil, err
		}
		events := convertico.ConvertSlices[proto.Message, *pb.EventDef](ms)

		return &pb.GetEventDefsReply{Events: events}, nil
	}

	func (s *ControlServer) AddEventDef(ctx context.Context, in *pb.AddEventDefRequest) (*pb.AddEventDefReply, error) {
		mydao, err := s.getDao(ctx, in.Event.MerchantId.Id, TBL_EVENT, &pb.EventDef{})
		if err != nil {
			return nil, err
		}
		var m proto.Message = in.Event
		ms, err := mydao.Add(ctx, m)
		if err != nil {
			return nil, err
		}
		in.Event.ID = strconv.FormatUint(ms, 10)
		return &pb.AddEventDefReply{Event: in.Event}, nil
	}
*/
