// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tableflip-test/pb"
	"time"

	"github.com/cloudflare/tableflip"
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
	pidFile := flag.String("pid-file", "./app.pid", "`Path` to pid file")
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	upg, err := tableflip.New(tableflip.Options{PIDFile: *pidFile})
	if err != nil {
		panic(err)
	}
	defer upg.Stop()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGUSR2)
		for s := range sig {
			if s == syscall.SIGUSR2 {
				upg.Upgrade()
			}
		}
	}()

	// create a wait group
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			default: // Client logic : greeting server.Before shutdown, client will Complete the processing
				// Contact the server and print out its response.
				log.Printf("start Greeting")
				time.Sleep(time.Second * 10)
				ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
				r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
				if err != nil {
					log.Fatalf("could not greet: %v", err)
				}
				log.Printf("end Greeting: %s", r.GetMessage())
			}
		}
	}()

	// wait exit
	// 通知父进程ready
	if err := upg.Ready(); err != nil {
		panic(err)
	}

	// graceful exit
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)

	select {
	case <-ch:
		log.Println("get Exit")
	case <-upg.Exit():
		done <- struct{}{}
		time.Sleep(20 * time.Second)
		log.Println("upg Exit")
	}
}
