package main

import (
	"flag"
	"github.com/songrenru/tag_service/proto"
	"github.com/songrenru/tag_service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

var grpcPort string
var httpPort string

func init() {
	flag.StringVar(&grpcPort, "grpc_port", "9091", "gRPC 启动端口号")
	flag.StringVar(&httpPort, "http_port", "8081", "http 启动端口号")
	flag.Parse()
}

func main() {
	errCh := make(chan error)

	go func() {
		err := RunHttpServer(httpPort)
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		err := RunGrpcServer(grpcPort)
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		log.Fatalf("run server err: %v", err)
	}
}

func RunHttpServer(port string) error {
	srv := http.NewServeMux()
	srv.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`pong`))
	})

	return http.ListenAndServe(":"+port, srv)
}

func RunGrpcServer(port string) error {
	srv := grpc.NewServer()
	proto.RegisterTagServiceServer(srv, server.NewTagServer())
	reflection.Register(srv)

	ln, _ := net.Listen("tcp", ":" + port)
	return srv.Serve(ln)
}
