package server

import (
	"context"

	"github.com/andrescosta/goico/pkg/converter"
	"github.com/andrescosta/goico/pkg/database"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/ctl/internal/dao"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

const (
	dbPath = ".\\db.db"
)

const (
	tblPackage     = "package"
	tblTenant      = "tenant"
	tblEnvironment = "environment"
)

const (
	environmentID = "environment_1"
)

type ControlServer struct {
	pb.UnimplementedControlServer

	daos map[string]*dao.DAO[proto.Message]

	db *database.Database

	bJobPackage *grpchelper.GrpcBroadcaster[*pb.UpdateToPackagesStrReply, proto.Message]

	bEnviroment *grpchelper.GrpcBroadcaster[*pb.UpdateToEnviromentStrReply, proto.Message]
}

func NewCotrolServer(ctx context.Context) (*ControlServer, error) {
	db, err := database.Open(dbPath)

	if err != nil {
		return nil, err
	}

	return &ControlServer{

		daos: make(map[string]*dao.DAO[proto.Message]),

		db: db,

		bJobPackage: grpchelper.StartBroadcaster[*pb.UpdateToPackagesStrReply, proto.Message](ctx),

		bEnviroment: grpchelper.StartBroadcaster[*pb.UpdateToEnviromentStrReply, proto.Message](ctx),
	}, nil
}

func (c *ControlServer) Close() error {
	c.bJobPackage.Stop()

	c.bEnviroment.Stop()

	return c.db.Close()
}

func (c *ControlServer) GetPackages(ctx context.Context, in *pb.GetJobPackagesRequest) (*pb.GetJobPackagesReply, error) {
	if in.ID != nil {
		p, err := c.getPackage(ctx, in.Tenant, *in.ID)

		if err != nil {
			return nil, err
		}
		if p != nil {
			return &pb.GetJobPackagesReply{Packages: []*pb.JobPackage{p}}, nil
		}
		return &pb.GetJobPackagesReply{}, nil
	}
	packages, err := c.getPackages(ctx, in.Tenant)
	if err != nil {
		return nil, err
	}
	return &pb.GetJobPackagesReply{Packages: packages}, nil
}

func (c *ControlServer) GetAllPackages(ctx context.Context, _ *pb.GetAllJobPackagesRequest) (*pb.GetAllJobPackagesReply, error) {
	ms, err := c.getTenants(ctx)
	if err != nil {
		return nil, err
	}

	packages := make([]*pb.JobPackage, 0)

	for _, me := range ms {
		mydao, err := c.getDao(ctx, me.ID, tblPackage, &pb.JobPackage{})

		if err != nil {
			return nil, err
		}

		ms, err := mydao.All(ctx)

		if err != nil {
			return nil, err
		}

		ps := converter.Slices[proto.Message, *pb.JobPackage](ms)

		packages = append(packages, ps...)
	}

	return &pb.GetAllJobPackagesReply{Packages: packages}, nil
}

func (c *ControlServer) AddPackage(ctx context.Context, in *pb.AddJobPackageRequest) (*pb.AddJobPackageReply, error) {
	mydao, err := c.getDao(ctx, in.Package.Tenant, tblPackage, &pb.JobPackage{})

	if err != nil {
		return nil, err
	}

	var m proto.Message = in.Package

	_, err = mydao.Add(ctx, m)

	if err != nil {
		return nil, err
	}

	c.broadcastAdd(in.Package)

	return &pb.AddJobPackageReply{Package: in.Package}, nil
}

func (c *ControlServer) UpdatePackage(ctx context.Context, in *pb.UpdateJobPackageRequest) (*pb.UpdateJobPackageReply, error) {
	mydao, err := c.getDao(ctx, in.Package.Tenant, tblPackage, &pb.JobPackage{})

	if err != nil {
		return nil, err
	}

	var m proto.Message = in.Package

	err = mydao.Update(ctx, m)

	if err != nil {
		return nil, err
	}

	c.broadcastUpdate(in.Package)

	return &pb.UpdateJobPackageReply{}, nil
}

func (c *ControlServer) DeletePackage(ctx context.Context, in *pb.DeleteJobPackageRequest) (*pb.DeleteJobPackageReply, error) {
	mydao, err := c.getDao(ctx, in.Package.Tenant, tblPackage, &pb.JobPackage{})

	if err != nil {
		return nil, err
	}

	err = mydao.Delete(ctx, in.Package.ID)

	if err != nil {
		return nil, err
	}

	c.broadcastDelete(in.Package)

	return &pb.DeleteJobPackageReply{}, nil
}

func (c *ControlServer) getPackages(ctx context.Context, tenant string) ([]*pb.JobPackage, error) {
	mydao, err := c.getDao(ctx, tenant, tblPackage, &pb.JobPackage{})

	if err != nil {
		return nil, err
	}

	ms, err := mydao.All(ctx)

	if err != nil {
		return nil, err
	}

	packages := converter.Slices[proto.Message, *pb.JobPackage](ms)

	return packages, nil
}

func (c *ControlServer) getPackage(ctx context.Context, tenant string, id string) (*pb.JobPackage, error) {
	mydao, err := c.getDao(ctx, tenant, tblPackage, &pb.JobPackage{})

	if err != nil {
		return nil, err
	}

	ms, err := mydao.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if ms != nil {
		return (*ms).(*pb.JobPackage), nil
	}
	return nil, nil
}

func (c *ControlServer) GetTenants(ctx context.Context, in *pb.GetTenantsRequest) (*pb.GetTenantsReply, error) {
	if in.ID != nil {
		t, err := c.getTenant(ctx, *in.ID)

		if err != nil {
			return nil, err
		}

		if t != nil {
			return &pb.GetTenantsReply{Tenants: []*pb.Tenant{t}}, nil
		}
		return &pb.GetTenantsReply{}, nil
	}
	ts, err := c.getTenants(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetTenantsReply{Tenants: ts}, nil
}

func (c *ControlServer) getTenants(ctx context.Context) ([]*pb.Tenant, error) {
	mydao, err := c.getDaoGen(ctx, tblTenant, &pb.Tenant{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.All(ctx)
	if err != nil {
		return nil, err
	}

	tenants := converter.Slices[proto.Message, *pb.Tenant](ms)
	return tenants, nil
}

func (c *ControlServer) getTenant(ctx context.Context, id string) (*pb.Tenant, error) {
	mydao, err := c.getDaoGen(ctx, tblTenant, &pb.Tenant{})

	if err != nil {
		return nil, err
	}

	ms, err := mydao.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if ms != nil {
		return (*ms).(*pb.Tenant), nil
	}
	return nil, nil
}

func (c *ControlServer) AddTenant(ctx context.Context, in *pb.AddTenantRequest) (*pb.AddTenantReply, error) {
	mydao, err := c.getDaoGen(ctx, tblTenant, &pb.Tenant{})
	if err != nil {
		return nil, err
	}

	var m proto.Message = in.Tenant
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	c.broadcastAdd(in.Tenant)
	return &pb.AddTenantReply{Tenant: in.Tenant}, nil
}

func (c *ControlServer) AddEnviroment(ctx context.Context, in *pb.AddEnviromentRequest) (*pb.AddEnviromentReply, error) {
	mydao, err := c.getDaoGen(ctx, tblEnvironment, &pb.Environment{})

	if err != nil {
		return nil, err
	}

	in.Environment.ID = environmentID

	var m proto.Message = in.Environment

	_, err = mydao.Add(ctx, m)

	if err != nil {
		return nil, err
	}

	c.broadcastAdd(in.Environment)

	return &pb.AddEnviromentReply{Environment: in.Environment}, nil
}

func (c *ControlServer) UpdateEnviroment(ctx context.Context, in *pb.UpdateEnviromentRequest) (*pb.UpdateEnviromentReply, error) {
	in.Environment.ID = environmentID

	mydao, err := c.getDaoGen(ctx, tblEnvironment, &pb.Environment{})

	if err != nil {
		return nil, err
	}

	var m proto.Message = in.Environment

	err = mydao.Update(ctx, m)

	if err != nil {
		return nil, err
	}

	c.broadcastUpdate(in.Environment)

	return &pb.UpdateEnviromentReply{}, nil
}

func (c *ControlServer) GetEnviroment(ctx context.Context, _ *pb.GetEnviromentRequest) (*pb.GetEnviromentReply, error) {
	mydao, err := c.getDaoGen(ctx, tblEnvironment, &pb.Environment{})
	if err != nil {
		return nil, err
	}

	ms, err := mydao.Get(ctx, environmentID)

	if err != nil {
		return nil, err
	}

	var environment *pb.Environment

	if ms != nil {
		environment = (*ms).(*pb.Environment)
	}

	return &pb.GetEnviromentReply{Environment: environment}, nil
}

func (c *ControlServer) UpdateToPackagesStr(_ *pb.UpdateToPackagesStrRequest, r pb.Control_UpdateToPackagesStrServer) error {
	return c.bJobPackage.RcvAndDispatchUpdates(r)
}

func (c *ControlServer) UpdateToEnviromentStr(_ *pb.UpdateToEnviromentStrRequest, r pb.Control_UpdateToEnviromentStrServer) error {
	return c.bEnviroment.RcvAndDispatchUpdates(r)
}

func (c *ControlServer) getDaoGen(ctx context.Context, entity string, message proto.Message) (*dao.DAO[proto.Message], error) {
	return c.getDao(ctx, entity, entity, message)
}

func (c *ControlServer) getDao(ctx context.Context, tenant string, entity string, message proto.Message) (*dao.DAO[proto.Message], error) {
	mydao, ok := c.daos[tenant]

	if !ok {
		var err error

		mydao, err = dao.NewDAO(ctx, c.db, tenant+"/"+entity,

			&ProtoMessageMarshaler{

				prototype: message,
			})

		if err != nil {
			return nil, err
		}

		c.daos[tenant] = mydao
	}

	return mydao, nil
}

func (c *ControlServer) broadcastAdd(m proto.Message) {
	c.broadcast(m, pb.UpdateType_New)
}

func (c *ControlServer) broadcastUpdate(m proto.Message) {
	c.broadcast(m, pb.UpdateType_Update)
}

func (c *ControlServer) broadcastDelete(m proto.Message) {
	c.broadcast(m, pb.UpdateType_Delete)
}

func (c *ControlServer) broadcast(m proto.Message, utype pb.UpdateType) {
	j, ok := m.(*pb.JobPackage)

	if ok {
		c.bJobPackage.Broadcast(j, utype)

		return
	}

	e, ok := m.(*pb.Environment)

	if ok {
		c.bEnviroment.Broadcast(e, utype)

		return
	}
}
