package core

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"rpcProject/common/api"
)

/*
	// Server端
	1. 在本地维护一个Call ID到函数指针的映射call_id_map，可以用std::map<std::string, std::function<>>
	2. 等待请求
	3. 得到一个请求后，将其数据包反序列化，得到Call ID
	4. 通过在call_id_map中查找，得到相应的函数指针
	5. 将lvalue和rvalue反序列化后，在本地调用Multiply函数，得到结果
	6. 将结果序列化后通过网络返回给Client
*/

// Server struct
type Server struct {
	addr  string
	funcs map[string]reflect.Value //1. 在本地维护一个Call ID到函数指针的映射call_id_map，可以用std::map<std::string, std::function<>>
}

// NewServer creates a new server
func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

//
func (s *Server) Run() {
	l, err := net.Listen("tcp", s.addr)

	if err != nil {
		log.Printf("listen on %s err: %v\n", s.addr, err)
		return
	}
	for {

		fmt.Printf("tcp server Run listener's network address: %v\n RemoteAddr %v", l.Addr())

		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
			continue
		}
		fmt.Println("====server Accept() client data ======")
		fmt.Printf("LocalAddr: %v\n | RemoteAddr %v", conn.LocalAddr(), conn.RemoteAddr())
		go func() {
			srvTransport := api.NewTransport(conn)

			for {
				// read request from client
				// 3. 得到一个请求后，将其数据包反序列化，得到Call ID
				req, err := srvTransport.Receive()
				if err != nil {
					if err != io.EOF {
						log.Printf("read err: %v\n", err)
					}
					return
				}
				// get method by name
				// 	4. 通过在call_id_map中查找，得到相应的函数指针
				f, ok := s.funcs[req.Name]
				if !ok { // if method requested does not exist
					e := fmt.Sprintf("func %s does not exist", req.Name)
					log.Println(e)
					if err = srvTransport.Send(api.Data{Name: req.Name, Err: e}); err != nil {
						log.Printf("transport write err: %v\n", err)
					}
					continue
				}
				log.Printf("func %s is called\n", req.Name)
				// unpackage request arguments
				inArgs := make([]reflect.Value, len(req.Args))
				for i := range req.Args {
					inArgs[i] = reflect.ValueOf(req.Args[i])
				}
				// invoke requested method
				// 	5. 将lvalue和rvalue反序列化后，在本地调用Multiply函数，得到结果
				out := f.Call(inArgs)
				// package response arguments (except error)
				outArgs := make([]interface{}, len(out)-1)
				for i := 0; i < len(out)-1; i++ {
					outArgs[i] = out[i].Interface()
				}
				// package error argument
				var e string
				if _, ok := out[len(out)-1].Interface().(error); !ok {
					e = ""
				} else {
					e = out[len(out)-1].Interface().(error).Error()
				}
				// send response to client
				// 		6. 将结果序列化后通过网络返回给Client
				err = srvTransport.Send(api.Data{Name: req.Name, Args: outArgs, Err: e})
				if err != nil {
					log.Printf("transport write err: %v\n", err)
				}
			}
		}()
	}
}

// Register a method via name
// 1. 在本地维护一个Call ID到函数指针的映射call_id_map，可以用std::map<std::string, std::function<>>
func (s *Server) Register(name string, f interface{}) {
	if _, ok := s.funcs[name]; ok {
		return
	}
	s.funcs[name] = reflect.ValueOf(f)
}
