package handler

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	proto "github.com/nebiros/poc-go-micro/client/proto/client"
)

type Client struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Client) Call(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Info("Received Client.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Client) Stream(ctx context.Context, req *proto.StreamingRequest, stream proto.Client_StreamStream) error {
	log.Infof("Received Client.Stream request with count: %d", req.Count)

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
func (e *Client) PingPong(ctx context.Context, stream proto.Client_PingPongStream) error {
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