package server

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Add  string
	C    chan string
	conn net.Conn

	server *Server
}

// NewUser 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Add:    userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听当前user channel 消息的goroutine
	go user.ListenMessage()

	return user
}

// Online 用户上线业务
func (u *User) Online() {
	// 用户上线，将用户降加入到onlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// 广播当前用户上线消息
	u.server.BroadCast(u, "已上线")
}

// Offline 用户下线业务
func (u *User) Offline() {
	// 用户下线，将用户从onlineMap中删除
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// 广播当前用户下线消息
	u.server.BroadCast(u, "下线")
}

func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// DoMessage 用户处理消息的业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//	查询当前在线用户列表
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Add + "]" + user.Name + ":" + "在线... \n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式 rename|张三
		newName := strings.Split(msg, "|")[1]

		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("当前用户名被使用\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("你已经更新用户名：" + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式 to|张三|信息内容

		//1 获取对方的用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("消息格式不正确，请使用 \"to|张三|你好啊\"格式，\n")
			return
		}
		// 2 根据用户名，得到对方User对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg("该用户名不存在 \n")
			return
		}
		// 3 获取消息内容 ，通过对方的User对象将信息内容发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("无消息内容，请重发\n")
			return
		}
		remoteUser.SendMsg(u.Name + " 对你说:" + content)

	} else {

		u.server.BroadCast(u, msg)
	}

}

//ListenMessage 监听当前User channel 的方法，一旦有消息，就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))

	}
}
