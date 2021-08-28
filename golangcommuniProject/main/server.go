package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//定义一个服务器
type Server struct {
	Ip string
	Port int
	//保存在线用户信息，因为全局所以需要加锁
	OnlineMap map[string]*User
	mapLock sync.RWMutex
	//消息广播的channel
	Message chan string
}

func (this *Server) Handler(connect net.Conn) {
	fmt.Println("建立用户连接成功 正在处理业务")

	user := NewUser(connect, this)

	user.Online()


	//监听用户是否活跃的channel
	isLive := make(chan bool)


	//接受用户发送的自定义消息
	go func() {

		buf := make([]byte, 4096)

		for {
			n, err := connect.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("buf read error", err)
				return
			}
			if n == 0 {
				user.Offline()
				return
			}

			inputmsg := string(buf[:n-1])//去掉\n

			//对用户输入的信息进行分门别类的处理
			user.DoMessage(inputmsg)

			//用户发了消息说明活跃
			isLive <- true
		}
	}()


	//超时强踢功能 监听isLive 来判定是否活跃
	go func() {
		for {
			select {

			case <- isLive:
				//当前用户活跃

			case <- time.After(10 * time.Second) ://10s后触发的定时器，每次执行都会重置
				user.SendMesage("超时下线")
				close(user.C)
				user.conn.Close()
				user.Offline()

				return
			}
		}
	}()

}

//服务器的启动方法
func (this *Server) Start() {

	//socket listen tcp 1
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}

	go this.ListenMessage()

	for {
		//accept tcp 2
		connect, err := listener.Accept()
		if err != nil {
			fmt.Println("connect err", err)
			continue
		}

		//do handler tcp 3
		go this.Handler(connect)
	}


	//close listen socket
	defer listener.Close()


}

//把准备广播上线的消息写入管道
func (this *Server) BoardCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

//监听Message，有消息就发送出去
func (this *Server) ListenMessage() {

	for {
		board := <- this.Message

		this.mapLock.Lock()

		for _, user := range this.OnlineMap {
			user.C <- board
		}

		this.mapLock.Unlock()
	}
}

//创建一个server的对外接口
func NewServer(ip string, port int,) *Server {
	server := &Server{Ip: ip, Port: port, Message: make(chan string), OnlineMap: make(map[string]*User)}
	return server
}

