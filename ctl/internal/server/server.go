package server

import (
	"context"

	"github.com/andrescosta/goico/pkg/convertico"
	"github.com/andrescosta/goico/pkg/database"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/ctl/internal/dao"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

const (
	dbPath         = ".\\db.db"
	tblEvent       = "event"
	tblQueue       = "queue"
	tblPackage     = "package"
	tblListener    = "listener"
	tblTenant      = "tenant"
	tblEnvironment = "environment"
	tblExecutor    = "executor"
	genMerchant    = "[Generic]"
	environmentId  = "enviroment_1"
)

type ControlServer struct {
	pb.UnimplementedControlServer
	daos        map[string]*dao.DAO[proto.Message]
	db          *database.Database
	bJobPackage *grpchelper.GrpcBroadcaster[*pb.UpdateToPackagesStrReply, proto.Message]
	bEnviroment *grpchelper.GrpcBroadcaster[*pb.UpdateToEnviromentStrReply, proto.Message]
}

func NewCotrolServer(ctx context.Context) (*ControlServer, error) {
	db, err := database.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &ControlServer{
		daos:        make(map[string]*dao.DAO[proto.Message]),
		db:          db,
		bJobPackage: grpchelper.StartBroadcaster[*pb.UpdateToPackagesStrReply, proto.Message](ctx),
		bEnviroment: grpchelper.StartBroadcaster[*pb.UpdateToEnviromentStrReply, proto.Message](ctx),
	}, nil
}

func (s *ControlServer) Close() error {
	s.bJobPackage.Stop()
	s.bEnviroment.Stop()
	return s.db.Close()
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
		mydao, err := s.getDao(ctx, me.ID, tblPackage, &pb.JobPackage{})
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
	mydao, err := s.getDao(ctx, in.Package.TenantId, tblPackage, &pb.JobPackage{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Package
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}

	s.broadcastAdd(in.Package)
	return &pb.AddJobPackageReply{Package: in.Package}, nil

}
func (s *ControlServer) getPackages(ctx context.Context, tenantId string) ([]*pb.JobPackage, error) {
	mydao, err := s.getDao(ctx, tenantId, tblPackage, &pb.JobPackage{})
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
	mydao, err := s.getDao(ctx, tenantId, tblPackage, &pb.JobPackage{})
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
	mydao, err := s.getDaoGen(ctx, tblTenant, &pb.Tenant{})
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
	mydao, err := s.getDaoGen(ctx, tblTenant, &pb.Tenant{})
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
	mydao, err := s.getDaoGen(ctx, tblTenant, &pb.Tenant{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Tenant
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	s.broadcastAdd(in.Tenant)
	return &pb.AddTenantReply{Tenant: in.Tenant}, nil
}

func (s *ControlServer) AddEnviroment(ctx context.Context, in *pb.AddEnviromentRequest) (*pb.AddEnviromentReply, error) {
	mydao, err := s.getDaoGen(ctx, tblEnvironment, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	in.Environment.ID = environmentId
	var m proto.Message = in.Environment
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	s.broadcastAdd(in.Environment)
	return &pb.AddEnviromentReply{Environment: in.Environment}, nil
}
func (s *ControlServer) UpdateEnviroment(ctx context.Context, in *pb.UpdateEnviromentRequest) (*pb.UpdateEnviromentReply, error) {
	in.Environment.ID = environmentId
	mydao, err := s.getDaoGen(ctx, tblEnvironment, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Environment
	err = mydao.Update(ctx, m)
	if err != nil {
		return nil, err
	}
	s.broadcastUpdate(in.Environment)

	return &pb.UpdateEnviromentReply{}, nil
}

func (s *ControlServer) GetEnviroment(ctx context.Context, in *pb.GetEnviromentRequest) (*pb.GetEnviromentReply, error) {
	mydao, err := s.getDaoGen(ctx, tblEnvironment, &pb.Environment{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.Get(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	var environment *pb.Environment = nil
	if ms != nil {
		environment = (*ms).(*pb.Environment)
	}
	return &pb.GetEnviromentReply{Environment: environment}, nil
}

func (s *ControlServer) UpdateToPackagesStr(in *pb.UpdateToPackagesStrRequest, r pb.Control_UpdateToPackagesStrServer) error {
	return s.bJobPackage.RcvAndDispatchUpdates(r)
}

func (s *ControlServer) UpdateToEnviromentStr(in *pb.UpdateToEnviromentStrRequest, r pb.Control_UpdateToEnviromentStrServer) error {
	return s.bEnviroment.RcvAndDispatchUpdates(r)
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
				prototype: message,
			})
		if err != nil {
			return nil, err
		}
		s.daos[tenant] = mydao
	}
	return mydao, nil
}

func (c *ControlServer) broadcastAdd(m proto.Message) {
	j, ok := m.(*pb.JobPackage)
	if ok {
		c.bJobPackage.Broadcast(j, pb.UpdateType_New)
		return
	}
	e, ok := m.(*pb.Environment)
	if ok {
		c.bEnviroment.Broadcast(e, pb.UpdateType_New)
		return
	}
}

func (c *ControlServer) broadcastUpdate(m proto.Message) {
	j, ok := m.(*pb.JobPackage)
	if ok {
		c.bJobPackage.Broadcast(j, pb.UpdateType_Update)
		return
	}
	e, ok := m.(*pb.Environment)
	if ok {
		c.bEnviroment.Broadcast(e, pb.UpdateType_Update)
		return
	}
}
