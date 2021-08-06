package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":2020")
	if err != nil {
		log.Fatalln("listen err: %v", err)
	}

	go broadCaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept occure error: %v \n", err)
			continue
		}

		go handleConn(conn)
	}
}

var (
	messageChannel = make(chan string, 8)
	leavingChannel = make(chan *User)
	enteringChannel = make(chan *User)
)

func broadCaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{}
		case user := <-leavingChannel:
			delete(users, user)
			close(user.MessageChannel)
		case msg := <-messageChannel:
			for user := range users {
				user.MessageChannel <- msg
			}
		}
	}
}

type User struct {
	ID int
	Addr string
	EnterAt time.Time
	MessageChannel chan string
}

func (u User) String() string {
	return u.Addr + ", UID:" + strconv.Itoa(u.ID) + ", Enter At:" +
		u.EnterAt.Format("2006-01-02 15:04:05+8000")
}

func handleConn(conn net.Conn)  {
	defer conn.Close()

	// 构建用户
	user := &User{
		ID:      GenUserId(),
		Addr:    conn.RemoteAddr().String(),
		EnterAt: time.Now(),
		MessageChannel: make(chan string, 8),
	}
	
	// 读写分离，写 chan
	go sendMessage(conn, user.MessageChannel)

	// 发送欢迎消息
	user.MessageChannel <- "welcome: " + user.String()
	messageChannel <- "user:" + strconv.Itoa(user.ID) + " has enter"

	// 注册用户，mutex
	enteringChannel <- user

	// 超时剔除用户
	var userActive = make(chan struct{})
	go func() {
		d := 5*time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	// 读写分离，读
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- strconv.Itoa(user.ID) + ":" + input.Text()
		
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}
	
	// 用户注销
	leavingChannel <- user
	messageChannel <- "user:" + strconv.Itoa(user.ID) + " has left"
}

func sendMessage(conn net.Conn, ch <-chan string)  {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

// 生成用户ID
var (
	globalId int
	idLocker sync.Mutex
)

func GenUserId() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalId++
	return globalId
}
