package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	connect net.Conn
	flag int
}


func (client *Client) Menu() bool {
	flag := 999

	fmt.Println("**************************************")
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Println("**************************************")

	_, err :=  fmt.Scanln(&flag)
	if err != nil || flag < 0 || flag > 3 {
		fmt.Println("输入模式错误", err)
		return false
	}else {
		client.flag = flag
		return true
	}
}


func (client *Client) PublicChat()  {
	fmt.Println("请输入公聊内容,exit 退出")

	var chatContent string = ""

	for chatContent != "exit" {
		_, err := fmt.Scanln(&chatContent)
		if err != nil {
			fmt.Println("输入的内容异常", err)
		}
		if chatContent == "exit" {
			break
		}

		_, err1 := client.connect.Write([]byte(chatContent))
		if err1 != nil {
			fmt.Println("输入的内容异常", err)
		}

		fmt.Println("请输入公聊内容,exit 退出")
		chatContent = ""
	}
}

func (client *Client) SelectOnlineUsers() {
	_, err := client.connect.Write([]byte("select all user"))
	if err != nil {
		fmt.Println("查询错误", err)
		return
	}

}

func (client *Client) PrivateChat() {

	//查询在线用户
	client.SelectOnlineUsers()
	//选择一个进入私聊
	fmt.Println("请选择要私聊的用户名，exit退出")
	var toName string = ""
	var chatContent string = ""

	for toName != "exit" {
		fmt.Scanln(&toName)
		if toName == "exit" {
			break
		}
		fmt.Println("请输入内容")
		fmt.Scanln(&chatContent)


		sendMsg := fmt.Sprintf("to|%s|%s", toName, chatContent)
		client.connect.Write([]byte(sendMsg))


		toName  = ""
		chatContent = ""
		client.SelectOnlineUsers()
		fmt.Println("请选择要私聊的用户名，exit退出")
	}



}


func (client *Client) UpdateName() bool {

	fmt.Println("请输入想更新的用户名")

	_, err := fmt.Scanln(client.Name)

	if err != nil {
		fmt.Println("更新用户名失败", err)
		return false
	}else {
		sendMsg := "rename|" + client.Name + "|\n"
		client.connect.Write([]byte(sendMsg))

		return true
	}
}



func (client *Client) Run()  {
	for client.flag != 0 {
		for client.Menu() == false {
		}

		switch client.flag {
		case 1:
			fmt.Println("公聊模式")
			client.PublicChat()
			break
		case 2:
			fmt.Println("私聊模式")
			client.PrivateChat()
			break
		case 3:
			fmt.Println("更新用户名")
			client.UpdateName()
			break
		}
	}


}




func NewClient(serverIp string, serverPort int) *Client {

	//创建对象
	client := &Client{}
	client.ServerIp = serverIp
	client.ServerPort = serverPort

	//连接server
	connet, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("客户端拨号异常", err)
	}

	client.connect = connet


	//返回对象
	return client

}

//处理server回应的消息
func (client *Client) DealRespounds() {
	io.Copy(os.Stdout, client.connect)
}



