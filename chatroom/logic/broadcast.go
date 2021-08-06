package logic

type broadcaster struct {
	users map[string]*User

	enteringChannel chan *User
	leavingChannel chan *User
	messageChannel chan *Message

	checkUserChannel chan string
	checkUserCanInChannel chan bool
}

var Broadcaster = &broadcaster{
	users:                 make(map[string]*User),
	enteringChannel:       make(chan *User),
	leavingChannel:        make(chan *User),
	messageChannel:        make(chan *Message, 10),
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

func (b *broadcaster) Start() {
	for {
		select { // select同一时间只会选取一个通道的返回，即串行执行，所以访问users没有并发问题
		case user := <-b.enteringChannel:
			b.messageChannel <- NewUserEnterMessage(user)
			b.users[user.Nickname] = user

			//b.sendUserList(user) // todo client判断消息类型，自行操作，不是更好？
		case user := <-b.leavingChannel:
			delete(b.users, user.Nickname)
			user.CloseMessageChannel()
			b.messageChannel <- NewUserLeaveMessage(user)
		case nickname := <- b.checkUserChannel:
			_, ok := b.users[nickname]
			b.checkUserCanInChannel <- !ok
		case msg := <-b.messageChannel:
			b.Broadcast(msg)
		}
	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	// todo 上锁?
	for _, user := range b.users {
		if user.UID == msg.User.UID {
			continue
		}
		user.MessageChannel <- msg
	}
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

func (b *broadcaster) GetUserList() []*User {
	// todo 上锁?
	list := make([]*User, 0, len(b.users))
	for _, user := range b.users {
		list = append(list, user)
	}
	return list
}

//func (b *broadcaster) sendUserList(user *User) {
//	list := b.GetUserList()
//	user.MessageChannel <-
//}

