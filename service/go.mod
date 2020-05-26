module github.com/nebiros/poc-go-micro/service

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v2 v2.7.0
	github.com/nebiros/poc-go-micro/broker/azqueue v0.0.0-20200526012004-5dd536c63ede
	google.golang.org/protobuf v1.23.0
)

replace github.com/nebiros/poc-go-micro/broker/azqueue => ../broker/azqueue
