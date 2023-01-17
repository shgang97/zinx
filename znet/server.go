package znet

import (
	"fmt"
	"net"
	"zinx/utils"
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
	// 当前 server 的消息管理模块，用来绑定 msgId 和 router
	MsgHandler ziface.IMsgHandler
}

func (server *Server) Start() {
	fmt.Printf("[Zinx] is starting...\n")
	fmt.Printf("ServerName: %s, IP: %s, Port: %d\n", server.Name, server.IP, server.Port)
	go func() {
		// 0. 开启消息队列及 worker 工作池
		server.MsgHandler.StartWorkerPool()
		// 1. 获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(server.IPVersion, fmt.Sprintf("%s:%d", server.IP, server.Port))
		if err != nil {
			fmt.Printf("resolve tcp addr error: %s\n", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(server.IPVersion, addr)
		if err != nil {
			fmt.Printf("listen %s error %s\n", server.IPVersion, err)
			return
		}
		fmt.Printf("start zinx server success, %s listening\n", server.Name)

		var cid uint32
		cid = 0
		// 3. 阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果客户端有连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("listener accept error: %s\n", err)
				continue
			}

			connection := NewConnection(conn, cid, server.MsgHandler)
			cid++
			// 启动当前的连接业务处理
			go connection.Start()
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

func (server *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	server.MsgHandler.AddRouter(msgId, router)
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.Config.Name,
		IPVersion:  "tcp4",
		IP:         utils.Config.Host,
		Port:       utils.Config.TcpPort,
		MsgHandler: NewMsgHandler(),
	}
	return s
}
