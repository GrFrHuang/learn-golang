package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"learn-golang/learn-grpc/protoFiles/pb"
	"github.com/micro/go-micro/server"
	"github.com/GrFrHuang/gox/log"
)

type GreeterHandler struct{}

var count = 0

func (g *GreeterHandler) Hello(ctx context.Context, req *grpc_file.HelloRequest, rsp *grpc_file.HelloResponse) error {
	count++
	fmt.Println(count)
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func (g *GreeterHandler) InsertGreeter(ctx context.Context, req *grpc_file.InsertRequest, rsp *grpc_file.InsertResponse) error {
	rsp.Greeter = &grpc_file.Greeter{Name: "hello " + req.Greeter.Name}
	log.Info(rsp)
	return nil
}

func (g *GreeterHandler) SelectGreeter(ctx context.Context, req *grpc_file.SelectRequest, rsp *grpc_file.SelectResponse) error {
	fmt.Println("run Select")
	return nil
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("micro.srv.greeter"),
		micro.Version("latest"),
	)
	// Init will parse the command line flags.
	service.Init(
		micro.AfterStop(func() error {
			log.Info("micro service be stopped  ...")
			return nil
		}),
		micro.AfterStart(func() error {
			log.Info("micro service start success ...")
			return nil
		}),
		// db.Init();
	)

	// Register handler
	grpc_file.RegisterGreeterServiceHandler(service.Server(), new(GreeterHandler), server.InternalHandler(true))

	//go excFunc()

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

//func excFunc() {
//	// create a new function
//	fnc := micro.NewFunction(
//		micro.Name("greeter"),
//	)
//
//	// init the command line
//	fnc.Init()
//
//	// register a handler
//	fnc.Handle(new(GreeterHandler))
//
//	// run the function
//	fnc.Run()
//}
