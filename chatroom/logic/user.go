package logic

import (
	"context"
	"errors"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"regexp"
	"sync/atomic"
	"time"
)

var System = &User{}

type User struct {
	UID int `json:"uid"`
	Nickname string `json:"nickname"`
	EnterAt time.Time `json:"enter_at"`
	Addr string `json:"addr"`

	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn
}

var globalUID uint32 = 0

func NewUser(nickname, addr string, conn *websocket.Conn) *User {
	user := &User{
		Nickname:       nickname,
		EnterAt:        time.Now(),
		Addr:           addr,
		MessageChannel: make(chan *Message, 32),
		conn:           conn,
	}

	user.UID = int(atomic.AddUint32(&globalUID, 1))

	return user
}

func (u *User) SendMessage(ctx context.Context) {
	for message := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, message)
	}
}

// CloseMessageChannel 避免 goroutine 泄露
func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err error
	)

	for {
		// 接受
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		// 分析content
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}
			return err
		}
		// broadcast
		sendMsg := NewMessage(u, receiveMsg["content"])
		sendMsg.Content = FilterSensitive(sendMsg.Content)
		reg := regexp.MustCompile(`@[^\s@]{3,20}`)
		sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)
		Broadcaster.Broadcast(sendMsg)
	}
}

func genToken() {}

func parseTokenAndValidate() {}

func macSha256() {}

func validateMAC() {}
