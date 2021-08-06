package main

import (
	"github.com/songrenru/tag_service/proto"
	"github.com/songrenru/tag_service/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

var port = "9091"

func main() {
	srv := grpc.NewServer()
	proto.RegisterTagServiceServer(srv, server.NewTagServer())

	ln, _ := net.Listen("tcp", ":" + port)
	err := srv.Serve(ln)
	if err != nil {
		log.Fatalf("srv.Serve err: %v", err)
	}
}
