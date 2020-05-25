package main

import (
	"context"
	"log"
	"time"

	"github.com/micro/go-micro/v2"

	proto "github.com/nebiros/poc-go-micro/pub/proto/pub"
)

func main() {
	service := micro.NewService(
		micro.Name("com.thriveglobal.service.poc.pub"),
		micro.Version("latest"),
	)

	service.Init()

	publisher := micro.NewPublisher("com.thriveglobal.service.poc", service.Client())

	for now := range time.Tick(time.Second) {
		if err := publisher.Publish(context.TODO(), &proto.Message{Say: now.String()}); err != nil {
			log.Fatalf("cannot publish: %s", err)
		}
	}
}
