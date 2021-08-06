package main

import (
	"fmt"
	_ "github.com/songrenru/chatroom/global"
	"github.com/songrenru/chatroom/server"
	"log"
	"net/http"
)

var (
	addr   = ":2022"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |

eason 的 Go 语言编程之旅 —— 一起用 Go 做项目：ChatRoom，start on：%s
`
)



func main() {
	fmt.Printf(banner+"\n", addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
