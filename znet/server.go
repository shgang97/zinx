package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

/*
@author: shg
@since: 2023/1/15
@desc: IServer 接口的实现，定义一个 Server 的服务器模块
*/

type Server struct {
	// 服务器名称
	Name string
	// 服务器监听的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的窗口
	Port int
}

func (server *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP: %s, Port: %d, is starting", server.IP, server.Port)
	go func() {
		// 1. 获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(server.IPVersion, fmt.Sprintf("%s:%d", server.IP, server.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(server.IPVersion, addr)
		if err != nil {
			fmt.Printf("listen %s error %s", server.IPVersion, err)
			return
		}
		fmt.Printf("start zinx server success, %s listening", server.Name)

		// 3. 阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("listener accept error: %s", err)
				continue
			}
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("read buf err: ", err)
						continue
					}
					fmt.Printf("recv client buf %s, cnt: %d\n", buf, cnt)
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err: ", err)
						continue
					}
				}
			}()
		}
	}()
}

func (server *Server) Stop() {
	// TODO
}

func (server *Server) Serve() {
	// 启动 server 的服务器功能
	server.Start()
	// TODO 做一些自动服务器之后的额外业务
	// 阻塞状态
	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
