package repo

import (
	"context"

	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	pb.UnimplementedRepoServer
	repo        *FileRepo
	bJobPackage *grpchelper.GrpcBroadcaster[*pb.UpdateToFileStrReply, proto.Message]
}

func NewServer(ctx context.Context, dir string) *Server {
	return &Server{
		repo:        NewFileRepo(dir),
		bJobPackage: grpchelper.StartBroadcaster[*pb.UpdateToFileStrReply, proto.Message](ctx),
	}
}

func (s *Server) AddFile(ctx context.Context, r *pb.AddFileRequest) (*pb.AddFileReply, error) {

	if err := s.repo.AddFile(r.TenantFile.TenantId, r.TenantFile.File.Name, int32(r.TenantFile.File.Type), r.TenantFile.File.Content); err != nil {
		return nil, err
	}
	s.bJobPackage.Broadcast(&pb.TenantFile{TenantId: r.TenantFile.TenantId, File: &pb.File{Name: r.TenantFile.File.Name, Content: r.TenantFile.File.Content}}, pb.UpdateType_New)
	return &pb.AddFileReply{}, nil
}

func (s *Server) GetFile(ctx context.Context, r *pb.GetFileRequest) (*pb.GetFileReply, error) {
	f, err := s.repo.File(r.TenantFile.TenantId, r.TenantFile.File.Name)
	if err != nil {
		return nil, err
	}
	m, err := s.repo.GetMetadataForFile(r.TenantFile.TenantId, r.TenantFile.File.Name)
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

func (s *Server) GetAllFileNames(ctx context.Context, r *pb.GetAllFileNamesRequest) (*pb.GetAllFileNamesReply, error) {
	f, err := s.repo.Files()
	if err != nil {
		return nil, err
	}
	return &pb.GetAllFileNamesReply{
		TenantFiles: f,
	}, nil
}

func (s *Server) UpdateToFileStr(in *pb.UpdateToFileStrRequest, r pb.Repo_UpdateToFileStrServer) error {
	return s.bJobPackage.RcvAndDispatchUpdates(r)
}
