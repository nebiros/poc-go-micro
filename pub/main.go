package main

import (
	"context"
	"time"

	"github.com/nebiros/poc-go-micro/broker/azqueue"

	"github.com/micro/go-micro/v2/client"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/cmd"

	"github.com/micro/go-micro/v2"

	log "github.com/micro/go-micro/v2/logger"
	proto "github.com/nebiros/poc-go-micro/server/proto/example"
)

var (
	storageAccountName string
	storageAccountKey  string
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
}

func run() {
	addCustomCmdFlags()

	if err := cmd.Init(); err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name("com.thriveglobal.service.poc.client.pub"),
		micro.Version("latest"),
		micro.Broker(azqueue.NewBroker(
			azqueue.StorageAccountName(storageAccountName),
			azqueue.StorageAccountKey(storageAccountKey),
		)),
	)

	service.Init()

	publisher := micro.NewPublisher(
		// this is the queue name at azure
		"sample-queue",
		service.Client(),
	)

	for now := range time.Tick(time.Second) {
		if err := publisher.Publish(context.TODO(), &proto.Message{Say: now.String()}, client.PublishContext(context.Background())); err != nil {
			log.Fatalf("cannot publish: %s", err)
		}
	}
}
