package main

import (
	"context"
	"flag"
	hellowworld "github.com/songrenru/grpc-demo/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "访问端口号")
	flag.Parse()
}

func main() {
	conn, _ := grpc.Dial(":" + port, grpc.WithInsecure())
	defer conn.Close()

	clnt := hellowworld.NewGreeterClient(conn)
	//SayHello(clnt)
	//SayList(clnt)
	//SayRecord(clnt)
	SayRoute(clnt)
}

func SayHello(client hellowworld.GreeterClient) error {
	resp, _ := client.SayHello(context.Background(), &hellowworld.HelloRequest{Name: "eason"})
	log.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}

func SayList(client hellowworld.GreeterClient) error {
	clnt, _ := client.SayList(context.Background(), &hellowworld.HelloRequest{Name: "eason"})
	for {
		resp, err := clnt.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("client.SayList resp: %s", resp.Message)
	}
	return nil
}

func SayRecord(client hellowworld.GreeterClient) error {
	stream, _ := client.SayRecord(context.Background())
	for i := 1; i <= 5; i++ {
		_ = stream.Send(&hellowworld.HelloRequest{Name: "eason" + strconv.Itoa(i)})
	}

	resp, _ := stream.CloseAndRecv()
	log.Printf("resp: %v", resp)

	return nil
}

func SayRoute(client hellowworld.GreeterClient) error {
	stream, _ := client.SayRoute(context.Background())
	for i := 1; i <= 5; i++ {
		_ = stream.Send(&hellowworld.HelloRequest{Name: "client.SayRoute:" + strconv.Itoa(i)})
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: %v", resp)
	}

	_ = stream.CloseSend()

	return nil
}
