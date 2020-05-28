package subscriber

import (
	"context"

	proto "github.com/nebiros/poc-go-micro/server/proto/example"

	log "github.com/micro/go-micro/v2/logger"
)

type Example struct{}

func (e *Example) Handle(ctx context.Context, msg *proto.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *proto.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
