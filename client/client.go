package client

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {

	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	// 返回对象
	return client

}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>> 链接服务器失败...")
		return
	}
	fmt.Println(">>>>>>>> 链接服务器成功 ...")

	// 启动服务
	select {}
}
