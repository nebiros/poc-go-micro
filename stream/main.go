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
		micro.Name("com.thriveglobal.service.poc.client.stream"),
		micro.Version("latest"),
	)

	service.Init()

	client := service.Client()

	req := client.NewRequest("com.thriveglobal.service.poc", "Example.Stream", &proto.StreamingRequest{})

	stream, err := client.Stream(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	if err := stream.Send(&proto.StreamingRequest{Count: int64(10)}); err != nil {
		log.Fatal(err)
	}

	for stream.Error() == nil {
		resp := &proto.StreamingResponse{}

		err := stream.Recv(resp)
		if err != nil {
			log.Error(err)
			break
		}

		log.Infof("Stream response: %d", resp.Count)
	}

	if stream.Error() != nil {
		log.Fatal(err)
	}

	if err := stream.Close(); err != nil {
		log.Fatal("cannot close stream service: %s", err)
	}
}
