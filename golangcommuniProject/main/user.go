package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {

	addr := conn.RemoteAddr().String()

	user := &User{Name: addr, Addr: addr, C: make(chan string), conn: conn, server: server}

	//每个用户都要监听
	go user.ListenMessage()

	return user
}

//监听User 中的 channel的方法,有消息就发送给客户端
func (this *User) ListenMessage() {
	for  {
		msg := <- this.C //读不到消息就会在这里阻塞

		this.conn.Write([]byte(msg + "\n"))
	}
}



func (this *User) Online() {
	//把用户加入到Onlilemap
	this.server.mapLock.Lock()

	this.server.OnlineMap[this.Name] = this

	this.server.mapLock.Lock()
	//把准备广播上线的消息写入管道
	this.server.BoardCast(this, "上线了")


}


func (this *User) Offline() {
	this.server.mapLock.Lock()

	delete(this.server.OnlineMap, this.Name)

	this.server.mapLock.Unlock()

	this.server.BoardCast(this, "下线了")


}

func (this *User) SendMesage(msg string) {
	this.conn.Write([]byte(msg))
}


func (this *User)  Send_p2p_Mesage(name string, msg string) {
	userobj, ok := this.server.OnlineMap[name]

	if ok {
		userobj.C <- msg
	}else {
		this.SendMesage("查无此人")
	}
}


//对用户输入的信息进行分门别类的处理
func (this *User) DoMessage(msg string) {

	if msg == "select all user" {
		this.server.mapLock.Lock()

		for username, _ := range this.server.OnlineMap {
			onlineMessage :=  username
			this.SendMesage(onlineMessage)
		}

		this.server.mapLock.Unlock()

	}else if msg[0:1] == "to" { //消息格式   to|name|content
		userName := strings.Split(msg, "|")[1]
		if userName == "" {
			this.SendMesage("to|name|content 是格式")
		}

		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMesage("无内容")
		}

		this.Send_p2p_Mesage(userName, content)

	}

	this.server.BoardCast(this, msg)

}
