// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/service/service.proto

package service

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Service service

func NewServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Service service

type Service interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Service_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Service_PingPongService, error)
}

type service struct {
	c    client.Client
	name string
}

func NewService(name string, c client.Client) Service {
	return &service{
		c:    c,
		name: name,
	}
}

func (c *service) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Service.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Service_StreamService, error) {
	req := c.c.NewRequest(c.name, "Service.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &serviceStream{stream}, nil
}

type Service_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type serviceStream struct {
	stream client.Stream
}

func (x *serviceStream) Close() error {
	return x.stream.Close()
}

func (x *serviceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *serviceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *serviceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *serviceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *service) PingPong(ctx context.Context, opts ...client.CallOption) (Service_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Service.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &servicePingPong{stream}, nil
}

type Service_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type servicePingPong struct {
	stream client.Stream
}

func (x *servicePingPong) Close() error {
	return x.stream.Close()
}

func (x *servicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *servicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *servicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *servicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *servicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Service service

type ServiceHandler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, Service_StreamStream) error
	PingPong(context.Context, Service_PingPongStream) error
}

func RegisterServiceHandler(s server.Server, hdlr ServiceHandler, opts ...server.HandlerOption) error {
	type service interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Service struct {
		service
	}
	h := &serviceHandler{hdlr}
	return s.Handle(s.NewHandler(&Service{h}, opts...))
}

type serviceHandler struct {
	ServiceHandler
}

func (h *serviceHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.ServiceHandler.Call(ctx, in, out)
}

func (h *serviceHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.ServiceHandler.Stream(ctx, m, &serviceStreamStream{stream})
}

type Service_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type serviceStreamStream struct {
	stream server.Stream
}

func (x *serviceStreamStream) Close() error {
	return x.stream.Close()
}

func (x *serviceStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *serviceStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *serviceStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *serviceStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *serviceHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.ServiceHandler.PingPong(ctx, &servicePingPongStream{stream})
}

type Service_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type servicePingPongStream struct {
	stream server.Stream
}

func (x *servicePingPongStream) Close() error {
	return x.stream.Close()
}

func (x *servicePingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *servicePingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *servicePingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *servicePingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *servicePingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
