package main

import (
	"fmt"
	"time"
	"learn-golang/learn-grpc/protoFiles/pb"
	"google.golang.org/grpc"
	"net"
)

type CounterServerImp struct {
}

//响应流式数据
func (c *CounterServerImp) Count(req *grpc_file.CountReq, stream grpc_file.Counter_CountServer) error {
	fmt.Printf("request from client. start: [%v]\n", req.GetStart())

	i := req.GetStart()
	for {
		i++
		stream.Send(&grpc_file.CountRes{
			Num: i,
		})
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

	return nil
}

func main() {
	server := grpc.NewServer()
	grpc_file.RegisterCounterServer(server, &CounterServerImp{})

	address, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	if err := server.Serve(address); err != nil {
		panic(err)
	}
}
