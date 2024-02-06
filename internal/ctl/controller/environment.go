package controller

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/ctl/data"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

const (
	tblEnvironment = "environment"
	environmentID  = "environment_1"
)

type EnvironmentController struct {
	ctx          context.Context
	daoCache     *data.DAOS
	bEnvironment *grpchelper.GrpcBroadcaster[*pb.UpdateToEnvironmentStrReply, proto.Message]
}

func NewEnvironmentController(ctx context.Context, db *database.Database) *EnvironmentController {
	return &EnvironmentController{
		ctx:          ctx,
		daoCache:     data.NewDAOS(db),
		bEnvironment: grpchelper.StartBroadcaster[*pb.UpdateToEnvironmentStrReply, proto.Message](ctx),
	}
}

func (c *EnvironmentController) Close() error {
	return c.bEnvironment.Stop()
}

func (c *EnvironmentController) AddEnvironment(in *pb.AddEnvironmentRequest) (*pb.AddEnvironmentReply, error) {
	mydao, err := c.daoCache.Generic(tblEnvironment, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	in.Environment.ID = environmentID
	var m proto.Message = in.Environment

	if err := mydao.Add(m); err != nil {
		return nil, err
	}
	c.broadcastAdd(in.Environment)
	return &pb.AddEnvironmentReply{Environment: in.Environment}, nil
}

func (c *EnvironmentController) UpdateEnvironment(in *pb.UpdateEnvironmentRequest) (*pb.Void, error) {
	in.Environment.ID = environmentID
	mydao, err := c.daoCache.Generic(tblEnvironment, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Environment
	err = mydao.Update(m)
	if err != nil {
		return nil, err
	}
	c.broadcastUpdate(in.Environment)
	return &pb.Void{}, nil
}

func (c *EnvironmentController) GetEnvironment() (*pb.EnvironmentReply, error) {
	mydao, err := c.daoCache.Generic(tblEnvironment, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	ms, err := mydao.Get(environmentID)
	if err != nil {
		return nil, err
	}
	var environment *pb.Environment
	if ms != nil {
		environment = (*ms).(*pb.Environment)
	}
	return &pb.EnvironmentReply{Environment: environment}, nil
}

func (c *EnvironmentController) UpdateToEnvironmentStr(_ *pb.Void, r pb.Control_UpdateToEnvironmentStrServer) error {
	return c.bEnvironment.RcvAndDispatchUpdates(c.ctx, r)
}

func (c *EnvironmentController) broadcastAdd(m *pb.Environment) {
	c.broadcast(m, pb.UpdateType_New)
}

func (c *EnvironmentController) broadcastUpdate(m *pb.Environment) {
	c.broadcast(m, pb.UpdateType_Update)
}

func (c *EnvironmentController) broadcast(m *pb.Environment, utype pb.UpdateType) {
	_ = c.bEnvironment.Broadcast(c.ctx, m, utype)
}
