// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"tableflip-test/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	for {
		time.Sleep(time.Second) // 每隔1秒请求一次
		// Contact the server and print out its response.
		log.Printf("start Greeting")
		ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("end Greeting: %s", r.GetMessage())
	}
}
