package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	proto "github.com/nebiros/poc-go-micro/client/proto/client"
)

type Client struct{}

func (e *Client) Handle(ctx context.Context, msg *proto.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *proto.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
