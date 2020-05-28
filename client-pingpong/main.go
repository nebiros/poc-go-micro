package main

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	proto "github.com/nebiros/poc-go-micro/server/proto/example"

	"github.com/micro/go-micro/v2"
)

func init() {
	log.Init(log.WithLevel(log.Level(log.TraceLevel)))
}

func main() {
	service := micro.NewService(
		micro.Name("com.thriveglobal.service.poc.client"),
		micro.Version("latest"),
	)

	service.Init()

	client := service.Client()

	req := client.NewRequest("com.thriveglobal.service.poc", "Example.PingPong", &proto.StreamingRequest{})

	stream, err := client.Stream(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		if err := stream.Send(&proto.Ping{Stroke: int64(i + 1)}); err != nil {
			log.Fatal(err)
		}

		resp := &proto.Pong{}

		if err := stream.Recv(resp); err != nil {
			log.Error(err)
			break
		}

		log.Infof("Sent ping %v got pong %v\n", i+1, resp.Stroke)
	}

	if stream.Error() != nil {
		log.Fatal(err)
	}

	if err := stream.Close(); err != nil {
		log.Fatal("cannot close ping pong service: %s", err)
	}
}
