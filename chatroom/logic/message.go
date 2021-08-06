package logic

import "time"

const (
	MsgTypeNormal = iota
	MsgTypeWelcome
	MsgTypeUserEnter
	MsgTypeUserLeave
	MsgTypeError
)

type Message struct {
	User *User `json:"user"`
	Type int `json:"type"`
	Content string `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	Ats []string `json:"ats"`

	//Users map[string]*User `json:"users"` // todo 为什么不是[]*User?
}

func NewMessage(user *User, content string) *Message {
	message := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
	return message
}

func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeWelcome,
		Content: user.Nickname + " 您好，欢迎加入聊天室！",
		MsgTime: time.Now(),
	}
}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserEnter,
		Content: user.Nickname + " 加入了聊天室",
		MsgTime: time.Now(),
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserLeave,
		Content: user.Nickname + " 离开了聊天室",
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(content string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}
