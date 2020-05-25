package subscriber

import (
	"context"

	proto "github.com/nebiros/poc-go-micro/service/proto/service"

	log "github.com/micro/go-micro/v2/logger"
)

type Service struct{}

func (e *Service) Handle(ctx context.Context, msg *proto.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *proto.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
