package controller

import (
	"context"

	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/internal/repo/provider"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

type Controller struct {
	repo        *provider.FileRepo
	bJobPackage *grpchelper.GrpcBroadcaster[*pb.UpdateToFileStrReply, proto.Message]
	ctx         context.Context
}

func New(ctx context.Context, dir string) *Controller {
	return &Controller{
		ctx:         ctx,
		repo:        provider.New(dir),
		bJobPackage: grpchelper.StartBroadcaster[*pb.UpdateToFileStrReply, proto.Message](ctx),
	}
}

func (s *Controller) AddFile(ctx context.Context, r *pb.AddFileRequest) (*pb.AddFileReply, error) {
	if err := s.repo.AddFile(r.TenantFile.Tenant, r.TenantFile.File.Name, int32(r.TenantFile.File.Type), r.TenantFile.File.Content); err != nil {
		return nil, err
	}
	s.bJobPackage.Broadcast(ctx, &pb.TenantFile{Tenant: r.TenantFile.Tenant, File: &pb.File{Name: r.TenantFile.File.Name, Content: r.TenantFile.File.Content}}, pb.UpdateType_New)
	return &pb.AddFileReply{}, nil
}

func (s *Controller) GetFile(_ context.Context, r *pb.GetFileRequest) (*pb.GetFileReply, error) {
	f, err := s.repo.File(r.TenantFile.Tenant, r.TenantFile.File.Name)
	if err != nil {
		return nil, err
	}
	m, err := s.repo.GetMetadataForFile(r.TenantFile.Tenant, r.TenantFile.File.Name)
	if err != nil {
		return nil, err
	}
	return &pb.GetFileReply{
		File: &pb.File{
			Content: f,
			Type:    pb.File_FileType(m.FileType),
		},
	}, nil
}

func (s *Controller) GetAllFileNames(_ context.Context, _ *pb.GetAllFileNamesRequest) (*pb.GetAllFileNamesReply, error) {
	f, err := s.repo.Files()
	if err != nil {
		return nil, err
	}
	return &pb.GetAllFileNamesReply{
		TenantFiles: f,
	}, nil
}

func (s *Controller) UpdateToFileStr(_ *pb.UpdateToFileStrRequest, r pb.Repo_UpdateToFileStrServer) error {
	return s.bJobPackage.RcvAndDispatchUpdates(s.ctx, r)
}
