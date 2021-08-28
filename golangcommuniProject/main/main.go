package main

import (
	"flag"
	"fmt"
)

var serverIP string
var serverPort int

//把全局变量放到flag中

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器ip（默认127.0.0.1）") //client -ip xxxx
	flag.IntVar(&serverPort, "port", 8888, "设置服务器port（默认8888）")
}


func main() {
	//命令行解析参数
	flag.Parse()

	//server := NewServer(serverIP, serverPort)
	//server.Start()
	
	client := NewClient(serverIP, serverPort)

	if client == nil {
		fmt.Println("连接服务器失败")
	}else {
		fmt.Println("连接服务器成功")
	}


	go client.DealRespounds()

	//启动客户端
	client.Run()




}
