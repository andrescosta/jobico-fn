// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: control.proto

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
	Control_GetMechants_FullMethodName      = "/Control/GetMechants"
	Control_AddMerchant_FullMethodName      = "/Control/AddMerchant"
	Control_AddPackage_FullMethodName       = "/Control/AddPackage"
	Control_GetAllPackages_FullMethodName   = "/Control/GetAllPackages"
	Control_GetPackages_FullMethodName      = "/Control/GetPackages"
	Control_UpdatePackage_FullMethodName    = "/Control/UpdatePackage"
	Control_GetEnviroment_FullMethodName    = "/Control/GetEnviroment"
	Control_AddEnviroment_FullMethodName    = "/Control/AddEnviroment"
	Control_UpdateEnviroment_FullMethodName = "/Control/UpdateEnviroment"
)

// ControlClient is the client API for Control service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ControlClient interface {
	GetMechants(ctx context.Context, in *GetMerchantsRequest, opts ...grpc.CallOption) (*GetMerchantsReply, error)
	AddMerchant(ctx context.Context, in *AddMerchantRequest, opts ...grpc.CallOption) (*AddMerchantReply, error)
	AddPackage(ctx context.Context, in *AddJobPackageRequest, opts ...grpc.CallOption) (*AddJobPackageReply, error)
	GetAllPackages(ctx context.Context, in *GetAllJobPackagesRequest, opts ...grpc.CallOption) (*GetAllJobPackagesReply, error)
	GetPackages(ctx context.Context, in *GetJobPackagesRequest, opts ...grpc.CallOption) (*GetJobPackagesReply, error)
	UpdatePackage(ctx context.Context, in *UpdateJobPackageRequest, opts ...grpc.CallOption) (*UpdateJobPackageReply, error)
	GetEnviroment(ctx context.Context, in *GetEnviromentRequest, opts ...grpc.CallOption) (*GetEnviromentReply, error)
	AddEnviroment(ctx context.Context, in *AddEnviromentRequest, opts ...grpc.CallOption) (*AddEnviromentReply, error)
	UpdateEnviroment(ctx context.Context, in *UpdateEnviromentRequest, opts ...grpc.CallOption) (*UpdateEnviromentReply, error)
}

type controlClient struct {
	cc grpc.ClientConnInterface
}

func NewControlClient(cc grpc.ClientConnInterface) ControlClient {
	return &controlClient{cc}
}

func (c *controlClient) GetMechants(ctx context.Context, in *GetMerchantsRequest, opts ...grpc.CallOption) (*GetMerchantsReply, error) {
	out := new(GetMerchantsReply)
	err := c.cc.Invoke(ctx, Control_GetMechants_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) AddMerchant(ctx context.Context, in *AddMerchantRequest, opts ...grpc.CallOption) (*AddMerchantReply, error) {
	out := new(AddMerchantReply)
	err := c.cc.Invoke(ctx, Control_AddMerchant_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) AddPackage(ctx context.Context, in *AddJobPackageRequest, opts ...grpc.CallOption) (*AddJobPackageReply, error) {
	out := new(AddJobPackageReply)
	err := c.cc.Invoke(ctx, Control_AddPackage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) GetAllPackages(ctx context.Context, in *GetAllJobPackagesRequest, opts ...grpc.CallOption) (*GetAllJobPackagesReply, error) {
	out := new(GetAllJobPackagesReply)
	err := c.cc.Invoke(ctx, Control_GetAllPackages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) GetPackages(ctx context.Context, in *GetJobPackagesRequest, opts ...grpc.CallOption) (*GetJobPackagesReply, error) {
	out := new(GetJobPackagesReply)
	err := c.cc.Invoke(ctx, Control_GetPackages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) UpdatePackage(ctx context.Context, in *UpdateJobPackageRequest, opts ...grpc.CallOption) (*UpdateJobPackageReply, error) {
	out := new(UpdateJobPackageReply)
	err := c.cc.Invoke(ctx, Control_UpdatePackage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) GetEnviroment(ctx context.Context, in *GetEnviromentRequest, opts ...grpc.CallOption) (*GetEnviromentReply, error) {
	out := new(GetEnviromentReply)
	err := c.cc.Invoke(ctx, Control_GetEnviroment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) AddEnviroment(ctx context.Context, in *AddEnviromentRequest, opts ...grpc.CallOption) (*AddEnviromentReply, error) {
	out := new(AddEnviromentReply)
	err := c.cc.Invoke(ctx, Control_AddEnviroment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlClient) UpdateEnviroment(ctx context.Context, in *UpdateEnviromentRequest, opts ...grpc.CallOption) (*UpdateEnviromentReply, error) {
	out := new(UpdateEnviromentReply)
	err := c.cc.Invoke(ctx, Control_UpdateEnviroment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ControlServer is the server API for Control service.
// All implementations must embed UnimplementedControlServer
// for forward compatibility
type ControlServer interface {
	GetMechants(context.Context, *GetMerchantsRequest) (*GetMerchantsReply, error)
	AddMerchant(context.Context, *AddMerchantRequest) (*AddMerchantReply, error)
	AddPackage(context.Context, *AddJobPackageRequest) (*AddJobPackageReply, error)
	GetAllPackages(context.Context, *GetAllJobPackagesRequest) (*GetAllJobPackagesReply, error)
	GetPackages(context.Context, *GetJobPackagesRequest) (*GetJobPackagesReply, error)
	UpdatePackage(context.Context, *UpdateJobPackageRequest) (*UpdateJobPackageReply, error)
	GetEnviroment(context.Context, *GetEnviromentRequest) (*GetEnviromentReply, error)
	AddEnviroment(context.Context, *AddEnviromentRequest) (*AddEnviromentReply, error)
	UpdateEnviroment(context.Context, *UpdateEnviromentRequest) (*UpdateEnviromentReply, error)
	mustEmbedUnimplementedControlServer()
}

// UnimplementedControlServer must be embedded to have forward compatible implementations.
type UnimplementedControlServer struct {
}

func (UnimplementedControlServer) GetMechants(context.Context, *GetMerchantsRequest) (*GetMerchantsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMechants not implemented")
}
func (UnimplementedControlServer) AddMerchant(context.Context, *AddMerchantRequest) (*AddMerchantReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMerchant not implemented")
}
func (UnimplementedControlServer) AddPackage(context.Context, *AddJobPackageRequest) (*AddJobPackageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPackage not implemented")
}
func (UnimplementedControlServer) GetAllPackages(context.Context, *GetAllJobPackagesRequest) (*GetAllJobPackagesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPackages not implemented")
}
func (UnimplementedControlServer) GetPackages(context.Context, *GetJobPackagesRequest) (*GetJobPackagesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPackages not implemented")
}
func (UnimplementedControlServer) UpdatePackage(context.Context, *UpdateJobPackageRequest) (*UpdateJobPackageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePackage not implemented")
}
func (UnimplementedControlServer) GetEnviroment(context.Context, *GetEnviromentRequest) (*GetEnviromentReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnviroment not implemented")
}
func (UnimplementedControlServer) AddEnviroment(context.Context, *AddEnviromentRequest) (*AddEnviromentReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddEnviroment not implemented")
}
func (UnimplementedControlServer) UpdateEnviroment(context.Context, *UpdateEnviromentRequest) (*UpdateEnviromentReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEnviroment not implemented")
}
func (UnimplementedControlServer) mustEmbedUnimplementedControlServer() {}

// UnsafeControlServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ControlServer will
// result in compilation errors.
type UnsafeControlServer interface {
	mustEmbedUnimplementedControlServer()
}

func RegisterControlServer(s grpc.ServiceRegistrar, srv ControlServer) {
	s.RegisterService(&Control_ServiceDesc, srv)
}

func _Control_GetMechants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMerchantsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).GetMechants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Control_GetMechants_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).GetMechants(ctx, req.(*GetMerchantsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_AddMerchant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMerchantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).AddMerchant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Control_AddMerchant_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).AddMerchant(ctx, req.(*AddMerchantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_AddPackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddJobPackageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).AddPackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Control_AddPackage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).AddPackage(ctx, req.(*AddJobPackageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_GetAllPackages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllJobPackagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).GetAllPackages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Control_GetAllPackages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).GetAllPackages(ctx, req.(*GetAllJobPackagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_GetPackages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJobPackagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).GetPackages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Control_GetPackages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).GetPackages(ctx, req.(*GetJobPackagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_UpdatePackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateJobPackageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).UpdatePackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Control_UpdatePackage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).UpdatePackage(ctx, req.(*UpdateJobPackageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_GetEnviroment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEnviromentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).GetEnviroment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Control_GetEnviroment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).GetEnviroment(ctx, req.(*GetEnviromentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_AddEnviroment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddEnviromentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).AddEnviroment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Control_AddEnviroment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).AddEnviroment(ctx, req.(*AddEnviromentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Control_UpdateEnviroment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEnviromentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServer).UpdateEnviroment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Control_UpdateEnviroment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServer).UpdateEnviroment(ctx, req.(*UpdateEnviromentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Control_ServiceDesc is the grpc.ServiceDesc for Control service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Control_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Control",
	HandlerType: (*ControlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMechants",
			Handler:    _Control_GetMechants_Handler,
		},
		{
			MethodName: "AddMerchant",
			Handler:    _Control_AddMerchant_Handler,
		},
		{
			MethodName: "AddPackage",
			Handler:    _Control_AddPackage_Handler,
		},
		{
			MethodName: "GetAllPackages",
			Handler:    _Control_GetAllPackages_Handler,
		},
		{
			MethodName: "GetPackages",
			Handler:    _Control_GetPackages_Handler,
		},
		{
			MethodName: "UpdatePackage",
			Handler:    _Control_UpdatePackage_Handler,
		},
		{
			MethodName: "GetEnviroment",
			Handler:    _Control_GetEnviroment_Handler,
		},
		{
			MethodName: "AddEnviroment",
			Handler:    _Control_AddEnviroment_Handler,
		},
		{
			MethodName: "UpdateEnviroment",
			Handler:    _Control_UpdateEnviroment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "control.proto",
}
