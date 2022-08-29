package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//	 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息channel的goroutime,一旦有消息就发送给全部的在线User
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message

		//将msg 发送给全部的在线User
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// BroadCast 广播消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Add + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	//  当前连接的业务

	fmt.Println("连接建立成功")

	user := NewUser(conn, s)

	// 用户上线，将用户降加入到onlineMap中
	user.Online()

	// 监听用户是否活跳的channel
	isLive := make(chan bool)

	// 广播当前用户上线消息
	s.BroadCast(user, "已上线")

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			read, err := conn.Read(buf)
			if read == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			// 提取用户的消息 （去掉 '\n' )
			msg := string(buf[:read-1])

			// 用户针对msg 进行消息处理
			user.DoMessage(msg)

			// 用户的任意消息，代表当前用户是一个活
			isLive <- true
		}
	}()
	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活动，应该重置定时器
			// 不做如何事情，是为了激活select 更新下面的定时器

		case <-time.After(time.Second * 300):
			// 已经超时
			//将当前User强制关闭

			user.SendMsg("你被踢了")

			//销毁用的资源
			close(user.C)

			//关闭连接
			conn.Close()

			//退出当前的handler
			return //runtime.Goexit()
		}
	}

}

// Start 启动服务器接口
func (s *Server) Start() {
	// socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// close listen socket
	defer listen.Close()

	// 启动监听Message的goroutime
	go s.ListenMessager()

	for {
		// accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen accept err:", err)
			continue
		}

		//do handler
		go s.Handler(conn)
	}
}
