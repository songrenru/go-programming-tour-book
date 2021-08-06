package server

import (
	"github.com/songrenru/chatroom/logic"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(writer http.ResponseWriter, request *http.Request) {
	// websocket.Accept 接受 client的websocket握手，并将连接升级到websocket
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项（通过第三个参数 AcceptOptions 设置）。
	// 换句话说，默认情况下，它不允许跨源请求。如果发生错误，Accept 将始终写入适当的响应
	conn, err := websocket.Accept(writer, request, nil)
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "内服务部错误")

	// 检查新用户
	nickname := request.FormValue("nickname")
	if l := len(nickname); l <= 2 || l > 20 {
		log.Println("nickname illegal:", err)
		wsjson.Write(request.Context(), conn, "非法昵称，昵称长度：3-20")
		conn.Close(websocket.StatusNormalClosure, "nickname illegal!")
		return
	}
	if !logic.Broadcaster.CanEnterRoom(nickname) { // 并发存在问题，CanEnterRoom()直接占用key？ 无并发问题，logic.Broadcaster.Start的设计决定的
		log.Println("nickname exists:", err)
		wsjson.Write(request.Context(), conn, "昵称已存在")
		conn.Close(websocket.StatusNormalClosure, "nickname exists!")
		return
	}
	user := logic.NewUser(nickname, request.RemoteAddr, conn)

	// 开启用户发送msg的goroutine
	go user.SendMessage(request.Context())

	// 发送欢迎消息
	user.MessageChannel <- logic.NewWelcomeMessage(user)

	// 广播新用户消息
	logic.Broadcaster.Broadcast(logic.NewUserEnterMessage(user))

	// 注册用户
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, " joins chat")

	// 接受消息
	err = user.ReceiveMessage(request.Context())

	// 用户离开
	logic.Broadcaster.UserLeaving(user)
	log.Println("user:", nickname, " leaving chat")

	// 后续清理
	if err != nil {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	} else {
		conn.Close(websocket.StatusNormalClosure, "")
	}

}
