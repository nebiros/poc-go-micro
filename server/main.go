package main

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/server"

	"github.com/micro/cli/v2"

	"github.com/nebiros/poc-go-micro/server/handler"
	"github.com/nebiros/poc-go-micro/server/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/nebiros/poc-go-micro/broker/azqueue"
	proto "github.com/nebiros/poc-go-micro/server/proto/example"
)

type SomeQueueProvider string

const (
	AWSSomeQueueProvider   SomeQueueProvider = "AWS"
	AzureSomeQueueProvider SomeQueueProvider = "Azure"
)

var (
	storageAccountName string
	storageAccountKey  string
	queueProvider      string
)

func init() {
	log.Init(log.WithLevel(log.Level(log.TraceLevel)))
}

func main() {
	run()
}

func addCustomCmdFlags() {
	app := cmd.App()
	app.Flags = append(app.Flags, &cli.StringFlag{
		Name:        "storage_account_name",
		Usage:       "Storage account name",
		EnvVars:     []string{"STORAGE_ACCOUNT_NAME"},
		Required:    true,
		Destination: &storageAccountName,
	})
	app.Flags = append(app.Flags, &cli.StringFlag{
		Name:        "storage_account_key",
		Usage:       "Storage account key",
		EnvVars:     []string{"STORAGE_ACCOUNT_KEY"},
		Required:    true,
		Destination: &storageAccountKey,
	})

	app.Flags = append(app.Flags, &cli.StringFlag{
		Name:        "queue_provider",
		Usage:       "Queue provider",
		EnvVars:     []string{"QUEUE_PROVIDER"},
		Required:    true,
		Destination: &queueProvider,
	})
}

func run() {
	addCustomCmdFlags()

	if err := cmd.Init(); err != nil {
		log.Fatal(err)
	}

	subscribeOptions := broker.NewSubscribeOptions(
		// we can use this option or the micro.RegisterSubscriber topic param
		broker.Queue("sample-queue"),
		azqueue.CreateQueue(true),
		azqueue.NumWorkers(16),
		azqueue.PoisonMessageDequeueThreshold(4),
	)

	var brokerOption micro.Option

	switch SomeQueueProvider(queueProvider) {
	case AWSSomeQueueProvider:
		brokerOption = micro.Broker(azqueue.NewBroker(
			azqueue.StorageAccountName(storageAccountName),
			azqueue.StorageAccountKey(storageAccountKey),
		))

	case AzureSomeQueueProvider:
		brokerOption = micro.Broker(azqueue.NewBroker(
			azqueue.StorageAccountName(storageAccountName),
			azqueue.StorageAccountKey(storageAccountKey),
		))
	}

	service := micro.NewService(
		micro.Name("com.thriveglobal.service.poc"),
		micro.Version("latest"),
		brokerOption,
	)

	service.Init()

	if err := proto.RegisterExampleHandler(service.Server(), new(handler.Example)); err != nil {
		log.Fatal(err)
	}

	if err := micro.RegisterSubscriber(
		// this is the queue name at azure
		"sample-queue",
		service.Server(),
		new(subscriber.Example),
		// we can pass context options this way
		server.SubscriberContext(subscribeOptions.Context),
	); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
