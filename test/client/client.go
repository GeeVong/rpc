package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"rpcProject/common/api"
	client "rpcProject/core"
	"rpcProject/test/public"
)

func main() {
	api.GetProcessPort()

	gob.Register(public.ResponseQueryUser{})

	//1. client调用client stub，这是一次本地过程调用
	//2. client stub将参数打包成一个消息，然后发送这个消息。打包过程也叫做 marshalling

	addr := "0.0.0.0:2333"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("dial error: %v\n", err)
	}
	cli := client.NewClient(conn)

	con := cli.GetConn()
	con.LocalAddr()
	fmt.Printf("client LocalAddr: %v\n RemoteAddr %v", con.LocalAddr(), con.RemoteAddr())

	var correctQuery func(int) (public.ResponseQueryUser, error)
	//var wrongQuery func(int) (public.ResponseQueryUser, error)

	//2. client所在的系统将消息发送给server cli.Call
	//server的的系统将收到的包传给server stub
	//server stub解包得到参数。 解包也被称作 unmarshalling
	//最后server stub调用服务过程. 返回结果按照相反的步骤传给client

	// 	1. 将这个调用映射为Call ID（name 就是标示）。这里假设用最简单的字符串当Call ID的方法

	for true {
		cli.Call("queryUser", &correctQuery)
		u, err := correctQuery(1)
		if err != nil {
			log.Printf("query error: %v\n", err)
		} else {
			log.Printf("query result: %v %v %v\n", u.Name, u.Age, u.Msg)
		}
		break
	}

	//u, err = correctQuery(2)
	//if err != nil {
	//	log.Printf("query error: %v\n", err)
	//} else {
	//	log.Printf("query result: %v %v %v\n", u.Name, u.Age, u.Msg)
	//}

	//cli.Call("queryUser", &wrongQuery)
	//u, err = wrongQuery(1)
	//if err != nil {
	//	log.Printf("query error: %v\n", err)
	//} else {
	//	log.Println(u)
	//}

	conn.Close()
}
