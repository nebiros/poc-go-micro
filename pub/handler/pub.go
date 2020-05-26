package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	proto "github.com/nebiros/poc-go-micro/pub/proto/pub"
)

type Pub struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Pub) Call(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Info("Received Pub.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Pub) Stream(ctx context.Context, req *proto.StreamingRequest, stream proto.Pub_StreamStream) error {
	log.Infof("Received Pub.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&proto.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Pub) PingPong(ctx context.Context, stream proto.Pub_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&proto.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
