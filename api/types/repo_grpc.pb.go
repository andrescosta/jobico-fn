// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: repo.proto

package types

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Repo_GetFile_FullMethodName         = "/Repo/GetFile"
	Repo_AddFile_FullMethodName         = "/Repo/AddFile"
	Repo_GetAllFileNames_FullMethodName = "/Repo/GetAllFileNames"
)

// RepoClient is the client API for Repo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RepoClient interface {
	GetFile(ctx context.Context, in *GetFileRequest, opts ...grpc.CallOption) (*GetFileReply, error)
	AddFile(ctx context.Context, in *AddFileRequest, opts ...grpc.CallOption) (*AddFileReply, error)
	GetAllFileNames(ctx context.Context, in *GetAllFileNamesRequest, opts ...grpc.CallOption) (*GetAllFileNamesReply, error)
}

type repoClient struct {
	cc grpc.ClientConnInterface
}

func NewRepoClient(cc grpc.ClientConnInterface) RepoClient {
	return &repoClient{cc}
}

func (c *repoClient) GetFile(ctx context.Context, in *GetFileRequest, opts ...grpc.CallOption) (*GetFileReply, error) {
	out := new(GetFileReply)
	err := c.cc.Invoke(ctx, Repo_GetFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repoClient) AddFile(ctx context.Context, in *AddFileRequest, opts ...grpc.CallOption) (*AddFileReply, error) {
	out := new(AddFileReply)
	err := c.cc.Invoke(ctx, Repo_AddFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repoClient) GetAllFileNames(ctx context.Context, in *GetAllFileNamesRequest, opts ...grpc.CallOption) (*GetAllFileNamesReply, error) {
	out := new(GetAllFileNamesReply)
	err := c.cc.Invoke(ctx, Repo_GetAllFileNames_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RepoServer is the server API for Repo service.
// All implementations must embed UnimplementedRepoServer
// for forward compatibility
type RepoServer interface {
	GetFile(context.Context, *GetFileRequest) (*GetFileReply, error)
	AddFile(context.Context, *AddFileRequest) (*AddFileReply, error)
	GetAllFileNames(context.Context, *GetAllFileNamesRequest) (*GetAllFileNamesReply, error)
	mustEmbedUnimplementedRepoServer()
}

// UnimplementedRepoServer must be embedded to have forward compatible implementations.
type UnimplementedRepoServer struct {
}

func (UnimplementedRepoServer) GetFile(context.Context, *GetFileRequest) (*GetFileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}
func (UnimplementedRepoServer) AddFile(context.Context, *AddFileRequest) (*AddFileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddFile not implemented")
}
func (UnimplementedRepoServer) GetAllFileNames(context.Context, *GetAllFileNamesRequest) (*GetAllFileNamesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllFileNames not implemented")
}
func (UnimplementedRepoServer) mustEmbedUnimplementedRepoServer() {}

// UnsafeRepoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RepoServer will
// result in compilation errors.
type UnsafeRepoServer interface {
	mustEmbedUnimplementedRepoServer()
}

func RegisterRepoServer(s grpc.ServiceRegistrar, srv RepoServer) {
	s.RegisterService(&Repo_ServiceDesc, srv)
}

func _Repo_GetFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepoServer).GetFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Repo_GetFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepoServer).GetFile(ctx, req.(*GetFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Repo_AddFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepoServer).AddFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Repo_AddFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepoServer).AddFile(ctx, req.(*AddFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Repo_GetAllFileNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllFileNamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepoServer).GetAllFileNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Repo_GetAllFileNames_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepoServer).GetAllFileNames(ctx, req.(*GetAllFileNamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Repo_ServiceDesc is the grpc.ServiceDesc for Repo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Repo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Repo",
	HandlerType: (*RepoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFile",
			Handler:    _Repo_GetFile_Handler,
		},
		{
			MethodName: "AddFile",
			Handler:    _Repo_AddFile_Handler,
		},
		{
			MethodName: "GetAllFileNames",
			Handler:    _Repo_GetAllFileNames_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "repo.proto",
}
