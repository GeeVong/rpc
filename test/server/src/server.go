package main

import (
	"encoding/gob"
	"errors"
	"rpcProject/common/api"
	server "rpcProject/core"
	"rpcProject/test/public"
	test "rpcProject/test/server/src/test"
)

func queryUser(uid int) (public.ResponseQueryUser, error) {
	db := make(map[int]public.User)
	db[0] = public.User{Name: "Jiahonzheng", Age: 70}
	db[1] = public.User{Name: "ChiuSinYing", Age: 75}
	db[2] = public.User{Name: "222222", Age: 75}
	if u, ok := db[uid]; ok {
		return public.ResponseQueryUser{User: u, Msg: "success"}, nil
	}
	return public.ResponseQueryUser{User: public.User{}, Msg: "fail"}, errors.New("uid is not in database")
}

func main() {

	api.GetProcessPort()

	test.Hello()
	gob.Register(public.ResponseQueryUser{})
	addr := "0.0.0.0:2333"
	srv := server.NewServer(addr)        // 绑定ip地址，已经进场id
	srv.Register("queryUser", queryUser) //

	go srv.Run() //  2。等待请求

	select {}
}
