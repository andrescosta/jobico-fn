package server

import (
	"context"

	pb "github.com/andrescosta/jobico/internal/api/types"
	"github.com/andrescosta/jobico/internal/repo/controller"
)

type Server struct {
	pb.UnimplementedRepoServer
	controller *controller.Controller
}

func New(ctx context.Context, dir string, o controller.Options) *Server {
	return &Server{
		controller: controller.New(ctx, dir, o),
	}
}

func (s *Server) Close() error {
	return s.controller.Close()
}

func (s *Server) AddFile(ctx context.Context, in *pb.AddFileRequest) (*pb.AddFileReply, error) {
	return s.controller.AddFile(ctx, in)
}

func (s *Server) File(ctx context.Context, in *pb.FileRequest) (*pb.FileReply, error) {
	return s.controller.File(ctx, in)
}

func (s *Server) AllFileNames(ctx context.Context, in *pb.Void) (*pb.AllFileNamesReply, error) {
	return s.controller.AllFileNames(ctx, in)
}

func (s *Server) UpdateToFileStr(in *pb.UpdateToFileStrRequest, ctl pb.Repo_UpdateToFileStrServer) error {
	return s.controller.UpdateToFileStr(in, ctl)
}
