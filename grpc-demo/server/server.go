package main

import (
	"context"
	"flag"
	hellowworld "github.com/songrenru/grpc-demo/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "启动端口号")
	flag.Parse()
}

type GreeterServer struct {}

func (s GreeterServer) SayHello(ctx context.Context, r *hellowworld.HelloRequest) (*hellowworld.HelloReply, error) {
	return &hellowworld.HelloReply{Message: "hello world," + r.Name}, nil
}

func (s GreeterServer) SayList(r *hellowworld.HelloRequest, srv hellowworld.Greeter_SayListServer) error {
	for i := 0; i < 6; i++ {
		_ = srv.Send(&hellowworld.HelloReply{Message: "hello.list," + r.Name})
	}

	return nil
}

func (s GreeterServer) SayRecord(stream hellowworld.Greeter_SayRecordServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&hellowworld.HelloReply{Message: "say.record"})
		}
		if err != nil {
			return err
		}

		log.Printf("req: %v", req)
	}

	return nil
}

func (s GreeterServer) SayRoute(stream hellowworld.Greeter_SayRouteServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		stream.Send(&hellowworld.HelloReply{Message: "server.SayRoute recv:" + req.Name})
		log.Printf("req: %v", req)
	}

	return nil
}

func main() {
	srv := grpc.NewServer()
	hellowworld.RegisterGreeterServer(srv, &GreeterServer{})
	ln, _ := net.Listen("tcp", ":" + port)
	srv.Serve(ln)
}
