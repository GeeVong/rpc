package core

import (
	"errors"
	"net"
	"reflect"
	"rpcProject/common/api"
)

/*
	// Client端
	//    int l_times_r = Call(ServerAddr, Multiply, lvalue, rvalue)

	1. 将这个调用映射为Call ID。这里假设用最简单的字符串当Call ID的方法
	2. 将Call ID，lvalue和rvalue序列化。可以直接将它们的值以二进制形式打包
	3. 把2中得到的数据包发送给ServerAddr，这需要使用网络传输层
	4. 等待服务器返回结果
	5. 如果服务器调用成功，那么就将结果反序列化，并赋给l_times_r
*/

// Client struct
type Client struct {
	conn net.Conn
}

// NewClient creates a new client
func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}

// Call transforms a function prototype into a function
func (c *Client) Call(name string, fptr interface{}) {
	container := reflect.ValueOf(fptr).Elem()

	// 000000fa2dff81030101044461746101ff8200010301044e616d65010c0001044172677301ff84000103457272010c0000001cff830201010e5b5d696e74657266616365207b7d01ff84000110000069ff82010971756572795573657201012872706350726f6a6563742f746573742f7075626c69632e526573706f6e7365517565727955736572ff8503010111526573706f6e736551756572795573657201ff8600010201045573657201ff880001034d7367010c00000023ff87030101045573657201ff8800010201044e616d65010c000103416765010400000020ff861c01010b4368697553696e59696e6701ff96000107737563636573730000
	f := func(req []reflect.Value) []reflect.Value {
		cliTransport := api.NewTransport(c.conn)

		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, container.Type().NumOut())
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}
		// package request arguments
		inArgs := make([]interface{}, 0, len(req))
		for i := range req {
			inArgs = append(inArgs, req[i].Interface())
		}
		// send request to server
		err := cliTransport.Send(api.Data{Name: name, Args: inArgs})
		if err != nil { // local network error or encode error
			return errorHandler(err)
		}
		// receive response from server
		rsp, err := cliTransport.Receive()
		if err != nil { // local network error or decode error
			return errorHandler(err)
		}
		if rsp.Err != "" { // remote server error
			return errorHandler(errors.New(rsp.Err))
		}

		if len(rsp.Args) == 0 {
			rsp.Args = make([]interface{}, container.Type().NumOut())
		}
		// unpackage response arguments
		numOut := container.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 { // unpackage arguments (except error)
				if rsp.Args[i] == nil { // if argument is nil (gob will ignore "Zero" in transmission), set "Zero" value
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(rsp.Args[i])
				}
			} else { // unpackage error argument
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}

		return outArgs
	}

	// 2. 将Call ID，lvalue和rvalue序列化。可以直接将它们的值以二进制形式打包
	container.Set(reflect.MakeFunc(container.Type(), f))
}

func (c *Client) GetConn() net.Conn {
	return c.conn
}
