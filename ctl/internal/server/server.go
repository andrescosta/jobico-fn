package server

import (
	"context"

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
	TBL_TENANT      = "tenant"
	TBL_ENVIRONMENT = "environment"
	TBL_EXECUTOR    = "executor"
	GEN_MERCHANT    = "[Generic]"
)

type ControlServer struct {
	pb.UnimplementedControlServer
	daos map[string]*dao.DAO[proto.Message]
	db   *database.Database
}

func NewCotrolServer() (*ControlServer, error) {
	db, err := database.Open(DB_PATH)
	if err != nil {
		return nil, err
	}
	return &ControlServer{
		daos: make(map[string]*dao.DAO[proto.Message]),
		db:   db,
	}, nil
}

func (s *ControlServer) Close() error {
	return s.db.Close()
}
func (s *ControlServer) getPackages(ctx context.Context, tenantId string) ([]*pb.JobPackage, error) {
	mydao, err := s.getDao(ctx, tenantId, TBL_PACKAGE, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.All(ctx)
	if err != nil {
		return nil, err
	}
	packages := convertico.SliceWithSlice[proto.Message, *pb.JobPackage](ms)

	return packages, nil
}
func (s *ControlServer) getPackage(ctx context.Context, tenantId string, id string) (*pb.JobPackage, error) {
	mydao, err := s.getDao(ctx, tenantId, TBL_PACKAGE, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if ms != nil {
		return (*ms).(*pb.JobPackage), nil
	} else {
		return nil, nil
	}
}

func (s *ControlServer) GetPackages(ctx context.Context, in *pb.GetJobPackagesRequest) (*pb.GetJobPackagesReply, error) {
	if in.ID != nil {
		p, err := s.getPackage(ctx, in.TenantId, *in.ID)
		if err != nil {
			return nil, err
		}
		if p != nil {
			return &pb.GetJobPackagesReply{Packages: []*pb.JobPackage{p}}, nil
		} else {
			return &pb.GetJobPackagesReply{}, nil
		}
	} else {
		packages, err := s.getPackages(ctx, in.TenantId)
		if err != nil {
			return nil, err
		}
		return &pb.GetJobPackagesReply{Packages: packages}, nil
	}
}

func (s *ControlServer) GetAllPackages(ctx context.Context, in *pb.GetAllJobPackagesRequest) (*pb.GetAllJobPackagesReply, error) {
	ms, err := s.getTenants(ctx)
	if err != nil {
		return nil, err
	}
	packages := make([]*pb.JobPackage, 0)
	for _, me := range ms {
		mydao, err := s.getDao(ctx, me.ID, TBL_PACKAGE, &pb.JobPackage{})
		if err != nil {
			return nil, err
		}

		ms, err := mydao.All(ctx)
		if err != nil {
			return nil, err
		}
		ps := convertico.SliceWithSlice[proto.Message, *pb.JobPackage](ms)
		packages = append(packages, ps...)
	}
	return &pb.GetAllJobPackagesReply{Packages: packages}, nil
}

func (s *ControlServer) AddPackage(ctx context.Context, in *pb.AddJobPackageRequest) (*pb.AddJobPackageReply, error) {
	mydao, err := s.getDao(ctx, in.Package.TenantId, TBL_PACKAGE, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Package
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	return &pb.AddJobPackageReply{Package: in.Package}, nil

}

func (s *ControlServer) GetTenants(ctx context.Context, in *pb.GetTenantsRequest) (*pb.GetTenantsReply, error) {
	if in.ID != nil {
		t, err := s.getTenant(ctx, *in.ID)
		if err != nil {
			return nil, err
		}
		if t != nil {
			return &pb.GetTenantsReply{Tenants: []*pb.Tenant{t}}, nil
		} else {
			return &pb.GetTenantsReply{}, nil
		}
	} else {
		m, err := s.getTenants(ctx)
		if err != nil {
			return nil, err
		}
		return &pb.GetTenantsReply{Tenants: m}, nil
	}
}

func (s *ControlServer) getTenants(ctx context.Context) ([]*pb.Tenant, error) {
	mydao, err := s.getDaoGen(ctx, TBL_TENANT, &pb.Tenant{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.All(ctx)
	if err != nil {
		return nil, err
	}
	tenants := convertico.SliceWithSlice[proto.Message, *pb.Tenant](ms)

	return tenants, nil
}

func (s *ControlServer) getTenant(ctx context.Context, id string) (*pb.Tenant, error) {
	mydao, err := s.getDaoGen(ctx, TBL_TENANT, &pb.Tenant{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if ms != nil {
		return (*ms).(*pb.Tenant), nil
	} else {
		return nil, nil
	}
}

func (s *ControlServer) AddTenant(ctx context.Context, in *pb.AddTenantRequest) (*pb.AddTenantReply, error) {
	mydao, err := s.getDaoGen(ctx, TBL_TENANT, &pb.Tenant{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Tenant
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	return &pb.AddTenantReply{Tenant: in.Tenant}, nil
}

func (s *ControlServer) AddEnviroment(ctx context.Context, in *pb.AddEnviromentRequest) (*pb.AddEnviromentReply, error) {
	mydao, err := s.getDaoGen(ctx, TBL_ENVIRONMENT, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	in.Environment.ID = "enviroment_1"
	var m proto.Message = in.Environment
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	return &pb.AddEnviromentReply{Environment: in.Environment}, nil
}
func (s *ControlServer) UpdateEnviroment(ctx context.Context, in *pb.UpdateEnviromentRequest) (*pb.UpdateEnviromentReply, error) {
	in.Environment.ID = "enviroment_1"
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

	ms, err := mydao.Get(ctx, "enviroment_1")
	if err != nil {
		return nil, err
	}
	var environment *pb.Environment = nil
	if ms != nil {
		environment = (*ms).(*pb.Environment)
	}
	return &pb.GetEnviromentReply{Environment: environment}, nil
}

func (s *ControlServer) getDaoGen(ctx context.Context, entity string, message proto.Message) (*dao.DAO[proto.Message], error) {
	return s.getDao(ctx, entity, entity, message)
}

func (s *ControlServer) getDao(ctx context.Context, tenant string, entity string, message proto.Message) (*dao.DAO[proto.Message], error) {
	mydao, ok := s.daos[tenant]
	if !ok {
		var err error
		mydao, err = dao.NewDAO(ctx, s.db, tenant+"/"+entity,
			&ProtoMessageMarshaler{
				newMessage: func() proto.Message {
					return message
				},
			})
		if err != nil {
			return nil, err
		}
		s.daos[tenant] = mydao
	}
	return mydao, nil
}

/*
	func (s *ControlServer) GetQueueDefs(ctx context.Context, in *pb.GetQueueDefsRequest) (*pb.GetQueueDefsReply, error) {
		mydao, err := s.getDao(ctx, in.TenantId.Id, TBL_QUEUE, &pb.QueueDef{})
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
		mydao, err := s.getDao(ctx, in.Queue.TenantId.Id, TBL_QUEUE, &pb.QueueDef{})
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
		mydao, err := s.getDao(ctx, in.TenantId.Id, TBL_EVENT, &pb.EventDef{})
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
		mydao, err := s.getDao(ctx, in.Event.TenantId.Id, TBL_EVENT, &pb.EventDef{})
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
