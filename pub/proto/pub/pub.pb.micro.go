// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/pub/pub.proto

package pub

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

// Api Endpoints for Pub service

func NewPubEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Pub service

type PubService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Pub_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Pub_PingPongService, error)
}

type pubService struct {
	c    client.Client
	name string
}

func NewPubService(name string, c client.Client) PubService {
	return &pubService{
		c:    c,
		name: name,
	}
}

func (c *pubService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Pub.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pubService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Pub_StreamService, error) {
	req := c.c.NewRequest(c.name, "Pub.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &pubServiceStream{stream}, nil
}

type Pub_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type pubServiceStream struct {
	stream client.Stream
}

func (x *pubServiceStream) Close() error {
	return x.stream.Close()
}

func (x *pubServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *pubServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *pubServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *pubServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *pubService) PingPong(ctx context.Context, opts ...client.CallOption) (Pub_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Pub.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &pubServicePingPong{stream}, nil
}

type Pub_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type pubServicePingPong struct {
	stream client.Stream
}

func (x *pubServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *pubServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *pubServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *pubServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *pubServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *pubServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Pub service

type PubHandler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, Pub_StreamStream) error
	PingPong(context.Context, Pub_PingPongStream) error
}

func RegisterPubHandler(s server.Server, hdlr PubHandler, opts ...server.HandlerOption) error {
	type pub interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Pub struct {
		pub
	}
	h := &pubHandler{hdlr}
	return s.Handle(s.NewHandler(&Pub{h}, opts...))
}

type pubHandler struct {
	PubHandler
}

func (h *pubHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.PubHandler.Call(ctx, in, out)
}

func (h *pubHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.PubHandler.Stream(ctx, m, &pubStreamStream{stream})
}

type Pub_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type pubStreamStream struct {
	stream server.Stream
}

func (x *pubStreamStream) Close() error {
	return x.stream.Close()
}

func (x *pubStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *pubStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *pubStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *pubStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *pubHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.PubHandler.PingPong(ctx, &pubPingPongStream{stream})
}

type Pub_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type pubPingPongStream struct {
	stream server.Stream
}

func (x *pubPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *pubPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *pubPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *pubPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *pubPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *pubPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}