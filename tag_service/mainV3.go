package main

import (
	"context"
	"flag"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/soheilhy/cmux"
	"github.com/songrenru/tag_service/internal/middleware"
	"github.com/songrenru/tag_service/proto"
	"github.com/songrenru/tag_service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

var port string

func init() {
	flag.StringVar(&port, "port", "9091", "启动端口号")
	flag.Parse()
}

func main() {
	l, err := RunTcpServer(port)
	if err != nil {
		log.Fatalf("Run TCP Server err: %v", err)
	}

	m := cmux.New(l)
	grpcLn := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpLn := m.Match(cmux.HTTP1Fast())

	grpcSrv := RunGrpcServer()
	httpSrv := RunHttpServer(port)
	go grpcSrv.Serve(grpcLn)
	go httpSrv.Serve(httpLn)

	err = m.Serve()
	if err != nil {
		log.Fatalf("Run Serve err: %v", err)
	}
}

func RunTcpServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

func RunHttpServer(port string) *http.Server {
	srv := http.NewServeMux()
	srv.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`pong`))
	})

	//prefix := "/swagger-ui/"
	//fileServer := http.FileServer(&assetfs.AssetFS{
	//	Asset:    swagger.Asset,
	//	AssetDir: swagger.AssetDir,
	//	Prefix:   "third_party/swagger-ui",
	//})
	//srv.Handle(prefix, http.StripPrefix(prefix, fileServer))

	return &http.Server{
		Addr: ":" + port,
		Handler: srv,
	}
}

func RunGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer( // 本质上还是拦截器里面再调拦截器
			middleware.AccessLog,
			middleware.ErrorLog,
			middleware.Recovery,
		)),
	}
	srv := grpc.NewServer(opts...)
	proto.RegisterTagServiceServer(srv, server.NewTagServer())
	reflection.Register(srv)

	return srv
}

func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好，拦截器")
	resp, err := handler(ctx, req)
	log.Println("再见，拦截器")
	return resp, err
}

func WorldInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好2，拦截器")
	resp, err := handler(ctx, req)
	log.Println("再见2，拦截器")
	return resp, err
}
