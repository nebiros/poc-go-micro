package main

import (
	"github.com/nebiros/poc-go-micro/service/handler"
	"github.com/nebiros/poc-go-micro/service/subscriber"

	"github.com/micro/go-micro/v2/server"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/nebiros/poc-go-micro/broker/azqueue"
	proto "github.com/nebiros/poc-go-micro/service/proto/service"
)

func main() {
	service := micro.NewService(
		micro.Name("com.thriveglobal.service.poc"),
		micro.Version("latest"),
		micro.Broker(azqueue.NewBroker(
			azqueue.StorageQueueName("sample-queue"),
			azqueue.StorageAccountName(""),
			azqueue.StorageAccountKey(""),
		)),
	)

	service.Init()

	if err := proto.RegisterServiceHandler(service.Server(), new(handler.Service)); err != nil {
		log.Fatal(err)
	}

	if err := micro.RegisterSubscriber(
		"com.thriveglobal.service.poc",
		service.Server(),
		new(subscriber.Service),
		server.SubscriberQueue("sample-queue"),
	); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
