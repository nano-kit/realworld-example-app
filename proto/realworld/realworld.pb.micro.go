// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/realworld/realworld.proto

package com_example_service_realworld

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

// Api Endpoints for Realworld service

func NewRealworldEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Realworld service

type RealworldService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Realworld_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Realworld_PingPongService, error)
}

type realworldService struct {
	c    client.Client
	name string
}

func NewRealworldService(name string, c client.Client) RealworldService {
	return &realworldService{
		c:    c,
		name: name,
	}
}

func (c *realworldService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Realworld.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realworldService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Realworld_StreamService, error) {
	req := c.c.NewRequest(c.name, "Realworld.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &realworldServiceStream{stream}, nil
}

type Realworld_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type realworldServiceStream struct {
	stream client.Stream
}

func (x *realworldServiceStream) Close() error {
	return x.stream.Close()
}

func (x *realworldServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *realworldServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *realworldServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *realworldServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *realworldService) PingPong(ctx context.Context, opts ...client.CallOption) (Realworld_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Realworld.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &realworldServicePingPong{stream}, nil
}

type Realworld_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type realworldServicePingPong struct {
	stream client.Stream
}

func (x *realworldServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *realworldServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *realworldServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *realworldServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *realworldServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *realworldServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Realworld service

type RealworldHandler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, Realworld_StreamStream) error
	PingPong(context.Context, Realworld_PingPongStream) error
}

func RegisterRealworldHandler(s server.Server, hdlr RealworldHandler, opts ...server.HandlerOption) error {
	type realworld interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Realworld struct {
		realworld
	}
	h := &realworldHandler{hdlr}
	return s.Handle(s.NewHandler(&Realworld{h}, opts...))
}

type realworldHandler struct {
	RealworldHandler
}

func (h *realworldHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.RealworldHandler.Call(ctx, in, out)
}

func (h *realworldHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.RealworldHandler.Stream(ctx, m, &realworldStreamStream{stream})
}

type Realworld_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type realworldStreamStream struct {
	stream server.Stream
}

func (x *realworldStreamStream) Close() error {
	return x.stream.Close()
}

func (x *realworldStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *realworldStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *realworldStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *realworldStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *realworldHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.RealworldHandler.PingPong(ctx, &realworldPingPongStream{stream})
}

type Realworld_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type realworldPingPongStream struct {
	stream server.Stream
}

func (x *realworldPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *realworldPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *realworldPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *realworldPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *realworldPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *realworldPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}