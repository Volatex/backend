// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/notification.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NotificationService_SendVerificationCode_FullMethodName = "/notification.NotificationService/SendVerificationCode"
	NotificationService_VerifyCode_FullMethodName           = "/notification.NotificationService/VerifyCode"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	SendVerificationCode(ctx context.Context, in *SendVerificationCodeRequest, opts ...grpc.CallOption) (*SendVerificationCodeResponse, error)
	VerifyCode(ctx context.Context, in *VerifyCodeRequest, opts ...grpc.CallOption) (*VerifyCodeResponse, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) SendVerificationCode(ctx context.Context, in *SendVerificationCodeRequest, opts ...grpc.CallOption) (*SendVerificationCodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendVerificationCodeResponse)
	err := c.cc.Invoke(ctx, NotificationService_SendVerificationCode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) VerifyCode(ctx context.Context, in *VerifyCodeRequest, opts ...grpc.CallOption) (*VerifyCodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VerifyCodeResponse)
	err := c.cc.Invoke(ctx, NotificationService_VerifyCode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility.
type NotificationServiceServer interface {
	SendVerificationCode(context.Context, *SendVerificationCodeRequest) (*SendVerificationCodeResponse, error)
	VerifyCode(context.Context, *VerifyCodeRequest) (*VerifyCodeResponse, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNotificationServiceServer struct{}

func (UnimplementedNotificationServiceServer) SendVerificationCode(context.Context, *SendVerificationCodeRequest) (*SendVerificationCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendVerificationCode not implemented")
}
func (UnimplementedNotificationServiceServer) VerifyCode(context.Context, *VerifyCodeRequest) (*VerifyCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyCode not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}
func (UnimplementedNotificationServiceServer) testEmbeddedByValue()                             {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	// If the following call pancis, it indicates UnimplementedNotificationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_SendVerificationCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendVerificationCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SendVerificationCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SendVerificationCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SendVerificationCode(ctx, req.(*SendVerificationCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_VerifyCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).VerifyCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_VerifyCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).VerifyCode(ctx, req.(*VerifyCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendVerificationCode",
			Handler:    _NotificationService_SendVerificationCode_Handler,
		},
		{
			MethodName: "VerifyCode",
			Handler:    _NotificationService_VerifyCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/notification.proto",
}
