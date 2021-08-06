package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/songrenru/tag_service/internal/middleware"
	"github.com/songrenru/tag_service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"log"
)

type Auth struct {
	AppKey string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func main() {
	auth := Auth{
		AppKey:    "go-programming-tour-book",
		AppSecret: "eason",
	}
	ctx := context.Background()
	md := metadata.New(map[string]string{"go": "programming", "author": "eason"})
	newCtx := metadata.NewOutgoingContext(ctx, md)

	target := "localhost:9091"
	opts := []grpc.DialOption{grpc.WithPerRPCCredentials(&auth)}
	conn, _ := GetClientConn(newCtx, target, opts)
	defer conn.Close()

	tagServiceClient := proto.NewTagServiceClient(conn)
	resp, err := tagServiceClient.GetTagList(newCtx, &proto.GetTagListRequest{})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err: %v", resp)
	}

	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.UnaryContextTimeout(),
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithMax(2),
				grpc_retry.WithCodes(
					codes.Unknown,
					codes.Internal,
					codes.DeadlineExceeded,
				),
			),
		),
	))
	opts = append(opts, grpc.WithStreamInterceptor(
		grpc_middleware.ChainStreamClient(
			middleware.StreamContextTimeout(),
		),
	))

	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}


