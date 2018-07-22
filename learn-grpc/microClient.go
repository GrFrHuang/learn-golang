package main

import (
	"context"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/client"
	"github.com/GrFrHuang/gox/log"
	"fmt"
	"learn-golang/learn-grpc/protoFiles/pb"
	"encoding/json"
	"github.com/micro/go-micro"
)

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("greeter.client"))
	service.Init()

	// Create new greeter client
	//greeter := grpc_file.NewGreeterService("micro.srv.greeter", service.Client())

	//var br = "你好吗"

	var request *grpc_file.HelloRequest
	//var a interface{}
	//bts2, _ := proto.Marshal(request)
	//fmt.Println("=====", string(bts2))
	hello := HelloW{
		Name: "world",
	}
	bts, _ := json.Marshal(hello)
	log.Info(string(bts))
	//var lt *grpc_file.HelloRequest

	json.Unmarshal(bts, &request)

	fmt.Println("=======", request)
	//response := grpc_file.HelloResponse{}

	//log.Info(proto.MarshalMessageSet(br))
	//bts, _ := json.Marshal(br)
	//request := json.RawMessage(bts)

	//var response json.RawMessage
	var response interface{}
	cli := *cmd.DefaultOptions().Client
	req := cli.NewRequest("micro.srv.greeter", "GreeterService.Hello", request, func(options *client.RequestOptions) {
		options.ContentType = "application/json"
	})
	err := cli.Call(context.TODO(), req, &response)

	// Call the greeter
	//rsp, err := greeter.InsertGreeter(context.TODO(), &grpc_file.InsertRequest{Greeter: &grpc_file.Greeter{Name: "GrFrHuang"}})
	if err != nil {
		log.Error(err)
	}
	bts2, _ := json.Marshal(response)
	fmt.Println(string(bts2))
	// Print response
	//fmt.Println(string(response))
	//fmt.Println(json.Marshal(string(response)))
}

type HelloW struct {
	Name string
}
