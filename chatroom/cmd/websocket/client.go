package main

import (
	"context"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://localhost:2021/ws", nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "内部错误")

	err = wsjson.Write(ctx, conn, "hello wesocket server")
	if err != nil {
		log.Println(err)
		return
	}

	var v interface{}
	err = wsjson.Read(ctx, conn, &v)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("接收到服务器响应：%v\n", v)

	conn.Close(websocket.StatusNormalClosure, "")
}
