package repo

import (
	"context"

	pb "github.com/andrescosta/jobico/api/types"
	"github.com/andrescosta/jobico/pkg/grpchelper"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	pb.UnimplementedRepoServer
	Repo        *FileRepo
	bJobPackage *grpchelper.GrpcBroadcaster[*pb.UpdateToFileStrReply, proto.Message]
}

func NewServer(ctx context.Context, dir string) *Server {
	return &Server{
		Repo:        &FileRepo{Dir: dir},
		bJobPackage: grpchelper.StartBroadcaster[*pb.UpdateToFileStrReply, proto.Message](ctx),
	}
}

func (s *Server) AddFile(ctx context.Context, r *pb.AddFileRequest) (*pb.AddFileReply, error) {
	err := s.Repo.AddFile(r.TenantId, r.Name, r.File)
	if err != nil {
		return nil, err
	}
	s.bJobPackage.Broadcast(&pb.UpdatedFile{TenantId: r.TenantId, Name: r.Name, File: r.File}, pb.UpdateType_New)
	return &pb.AddFileReply{}, nil
}

func (s *Server) GetFile(ctx context.Context, r *pb.GetFileRequest) (*pb.GetFileReply, error) {
	f, err := s.Repo.File(r.TenantId, r.Name)
	if err != nil {
		return nil, err
	}
	return &pb.GetFileReply{
		File: f,
	}, nil
}

func (s *Server) GetAllFileNames(ctx context.Context, r *pb.GetAllFileNamesRequest) (*pb.GetAllFileNamesReply, error) {
	f, err := s.Repo.Files()
	if err != nil {
		return nil, err
	}
	return &pb.GetAllFileNamesReply{
		Files: f,
	}, nil
}

func (s *Server) UpdateToFileStr(in *pb.UpdateToFileStrRequest, r pb.Repo_UpdateToFileStrServer) error {
	return s.bJobPackage.RcvAndDispatchUpdates(r)
}
