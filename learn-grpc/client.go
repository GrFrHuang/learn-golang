package main

import (
	"strconv"
	"os"
	"google.golang.org/grpc"
	"learn-golang/learn-grpc/protoFiles"
	"golang.org/x/net/context"
	"fmt"
)

func main() {
	start, _ := strconv.ParseInt(os.Args[1], 10, 64)

	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		fmt.Errorf("dial failed. err: [%v]\n", err)
		return
	}
	client := grpc_file.NewCounterClient(conn)

	stream, err := client.Count(context.Background(), &grpc_file.CountReq{
		Start: start,
	})
	if err != nil {
		fmt.Errorf("count failed. err: [%v]\n", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			fmt.Errorf("client count failed. err: [%v]", err)
			return
		}

		fmt.Printf("server count: %v\n", res.GetNum())
	}
}
