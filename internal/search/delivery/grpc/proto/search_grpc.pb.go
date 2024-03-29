// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: search.proto

package search

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
	SearchEngine_Get_FullMethodName = "/search.SearchEngine/Get"
)

// SearchEngineClient is the client API for SearchEngine service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchEngineClient interface {
	Get(ctx context.Context, in *Query, opts ...grpc.CallOption) (*QueryResult, error)
}

type searchEngineClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchEngineClient(cc grpc.ClientConnInterface) SearchEngineClient {
	return &searchEngineClient{cc}
}

func (c *searchEngineClient) Get(ctx context.Context, in *Query, opts ...grpc.CallOption) (*QueryResult, error) {
	out := new(QueryResult)
	err := c.cc.Invoke(ctx, SearchEngine_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchEngineServer is the server API for SearchEngine service.
// All implementations must embed UnimplementedSearchEngineServer
// for forward compatibility
type SearchEngineServer interface {
	Get(context.Context, *Query) (*QueryResult, error)
	mustEmbedUnimplementedSearchEngineServer()
}

// UnimplementedSearchEngineServer must be embedded to have forward compatible implementations.
type UnimplementedSearchEngineServer struct {
}

func (UnimplementedSearchEngineServer) Get(context.Context, *Query) (*QueryResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedSearchEngineServer) mustEmbedUnimplementedSearchEngineServer() {}

// UnsafeSearchEngineServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchEngineServer will
// result in compilation errors.
type UnsafeSearchEngineServer interface {
	mustEmbedUnimplementedSearchEngineServer()
}

func RegisterSearchEngineServer(s grpc.ServiceRegistrar, srv SearchEngineServer) {
	s.RegisterService(&SearchEngine_ServiceDesc, srv)
}

func _SearchEngine_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Query)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchEngineServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchEngine_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchEngineServer).Get(ctx, req.(*Query))
	}
	return interceptor(ctx, in, info, handler)
}

// SearchEngine_ServiceDesc is the grpc.ServiceDesc for SearchEngine service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SearchEngine_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "search.SearchEngine",
	HandlerType: (*SearchEngineServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _SearchEngine_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search.proto",
}
