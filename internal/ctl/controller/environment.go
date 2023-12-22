package controller

import (
	"context"

	"github.com/andrescosta/goico/pkg/database"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/ctl/dao"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

const (
	tblEnvironment = "environment"
	environmentID  = "environment_1"
)

type EnvironmentController struct {
	daoCache    *dao.Cache
	bEnviroment *grpchelper.GrpcBroadcaster[*pb.UpdateToEnviromentStrReply, proto.Message]
}

func NewEnvironmentController(ctx context.Context, db *database.Database) *EnvironmentController {
	return &EnvironmentController{
		daoCache:    dao.NewCache(db),
		bEnviroment: grpchelper.StartBroadcaster[*pb.UpdateToEnviromentStrReply, proto.Message](ctx),
	}
}
func (c *EnvironmentController) Close() {
	c.bEnviroment.Stop()
}

func (c *EnvironmentController) AddEnviroment(ctx context.Context, in *pb.AddEnviromentRequest) (*pb.AddEnviromentReply, error) {
	mydao, err := c.daoCache.GetGeneric(ctx, tblEnvironment, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	in.Environment.ID = environmentID
	var m proto.Message = in.Environment
	_, err = mydao.Add(ctx, m)
	if err != nil {
		return nil, err
	}
	c.broadcastAdd(ctx, in.Environment)
	return &pb.AddEnviromentReply{Environment: in.Environment}, nil
}
func (c *EnvironmentController) UpdateEnviroment(ctx context.Context, in *pb.UpdateEnviromentRequest) (*pb.UpdateEnviromentReply, error) {
	in.Environment.ID = environmentID
	mydao, err := c.daoCache.GetGeneric(ctx, tblEnvironment, &pb.Environment{})
	if err != nil {
		return nil, err
	}
	var m proto.Message = in.Environment
	err = mydao.Update(ctx, m)
	if err != nil {
		return nil, err
	}
	c.broadcastUpdate(ctx, in.Environment)
	return &pb.UpdateEnviromentReply{}, nil
}
func (c *EnvironmentController) GetEnviroment(ctx context.Context, _ *pb.GetEnviromentRequest) (*pb.GetEnviromentReply, error) {
	mydao, err := c.daoCache.GetGeneric(ctx, tblEnvironment, &pb.Environment{})
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
func (c *EnvironmentController) UpdateToEnviromentStr(_ *pb.UpdateToEnviromentStrRequest, r pb.Control_UpdateToEnviromentStrServer) error {
	return c.bEnviroment.RcvAndDispatchUpdates(r)
}
func (c *EnvironmentController) broadcastAdd(ctx context.Context, m *pb.Environment) {
	c.broadcast(ctx, m, pb.UpdateType_New)
}
func (c *EnvironmentController) broadcastUpdate(ctx context.Context, m *pb.Environment) {
	c.broadcast(ctx, m, pb.UpdateType_Update)
}
func (c *EnvironmentController) broadcast(ctx context.Context, m *pb.Environment, utype pb.UpdateType) {
	c.bEnviroment.Broadcast(ctx, m, utype)
}
