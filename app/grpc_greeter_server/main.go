// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"tableflip-test/pb"
	"time"

	"github.com/cloudflare/tableflip"

	"google.golang.org/grpc"
)

var (
	x    int64
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	i := atomic.AddInt64(&x, 1)
	log.Printf("Received: %v", in.GetName())
	time.Sleep(10 * time.Second)
	return &pb.HelloReply{Message: fmt.Sprintf("Req:%d, Hello %s", i, in.GetName())}, nil
}

func main() {
	pidFile := flag.String("pid-file", "./app.pid", "`Path` to pid file")
	flag.Parse()

	log.SetPrefix(fmt.Sprintf("%d ", os.Getpid()))
	upg, err := tableflip.New(tableflip.Options{
		PIDFile: *pidFile,
	})
	if err != nil {
		panic(err)
	}
	defer upg.Stop()
	// Do an upgrade on SIGHUP
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGUSR2)
		for range sig {
			log.Printf("reciece sigusr2")
			err := upg.Upgrade()
			if err != nil {
				log.Println("Upgrade failed:", err)
			}
		}
	}()

	ln, err := upg.Fds.Listen("tcp", "localhost:"+strconv.Itoa(*port))
	if err != nil {
		panic(err)
	}
	// prepare rpc server
	s := grpc.NewServer()                  // create rpc server
	pb.RegisterGreeterServer(s, &server{}) // register service
	go func() {
		err := s.Serve(ln)
		log.Printf("server listening at %v", ln.Addr())
		if err != http.ErrServerClosed {
			log.Printf("failed to serve: %v", err)
		}
	}()

	log.Printf("ready")
	if err := upg.Ready(); err != nil {
		panic(err)
	}
	<-upg.Exit()

	// Make sure to set a deadline on exiting the process
	// after upg.Exit() is closed. No new upgrades can be
	// performed if the parent doesn't exit.
	time.AfterFunc(30*time.Second, func() {
		log.Println("Graceful shutdown timed out")
		os.Exit(1)
	})

	// Wait for connections to drain.
	log.Println("stops the gRPC server gracefully")
	s.GracefulStop()
}
