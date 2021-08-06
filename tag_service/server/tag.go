package server

import (
	"context"
	"encoding/json"
	"github.com/songrenru/tag_service/pkg/bapi"
	"github.com/songrenru/tag_service/pkg/errcode"
	"github.com/songrenru/tag_service/proto"
	"google.golang.org/grpc/metadata"
	"log"
)

type TagServer struct {
	auth *Auth
}

type Auth struct {}

func (a *Auth) GetAppKey() string {
	return "go-programming-tour-book"
}

func (a *Auth) GetAppSecret() string {
	return "eason"
}

func (a *Auth) Check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	var appKey, appSecret string
	if val := md.Get("app_key"); len(val) > 0 {
		appKey = val[0]
	}
	if val := md.Get("app_secret"); len(val) > 0 {
		appSecret = val[0]
	}

	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}

	return nil
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, request *proto.GetTagListRequest) (*proto.GetTagListReply, error) {
	//panic("panic测试")
	//time.Sleep(time.Second*11)
	if err := t.auth.Check(ctx); err != nil {
		return nil, err
	}

	api := bapi.NewAPi("http://127.0.0.1:8080")
	body, err := api.GetTagList(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	tagList := proto.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, err
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("metadata: %v", md)
	}


	return &tagList, nil
}
