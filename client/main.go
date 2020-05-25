package main

import (
	"context"
	"fmt"
	"log"

	proto "github.com/nebiros/poc-go-micro/client/proto/client"

	"github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService(
		micro.Name("com.thriveglobal.service.poc.client"),
		micro.Version("latest"),
	)

	service.Init()

	client := proto.NewClientService("com.thriveglobal.service.poc", service.Client())

	resp, err := client.Call(context.TODO(), &proto.Request{Name: "Bill 4"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Msg)
}
