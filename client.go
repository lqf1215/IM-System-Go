package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前client的模式
}

func NewClient(serverIp string, serverPort int) *Client {

	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>> 请输入合法范围内容的数组<<<<<<")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		// 根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式选择...")
			break
		case 2:
			// 私聊模式
			fmt.Println("私聊模式...")

			break
		case 3:
			// 更新用户名
			//fmt.Println("更新用户名选择...")
			client.UpdateName()

			break
		}
	}
}

// 处理server回应的信息，直接显示到标准输出即可
func (client *Client) DealResponse() {
	// 一旦 client.conn 有数据，就直接copy 到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)

	//for {
	//	buf := make()
	//}
}

func (client *Client) UpdateName() bool {

	fmt.Println(">>>>>> 请输入用户名：")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

var serverIp string

var serverPort int

// ./client -ip 127.0.0.1

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认是8888）")
}

func main() {

	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>> 链接服务器失败...")
		return
	}

	// 单独启动一个goroutine去处理server的回执消息
	go client.DealResponse()

	fmt.Println(">>>>>>>> 链接服务器成功 ...")

	// 启动客户端的服务业务
	//select {}
	client.Run()
}
