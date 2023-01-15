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
@desc: //TODO
*/

type Connection struct {
	// 当前连接 socket tcp 套接字
	Conn *net.TCPConn
	// 连接的 ID
	ConnID uint32
	// 当前的连接状态
	IsClosed bool
	// 告知当前连接已经退出/停止的 channel
	ExitChan chan bool
	// 该连接处理的方法 router
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		IsClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
	return c
}

// StartReader 读业务方法
func (c *Connection) StartReader() {
	fmt.Printf("Conn start reading..., connID = %d\n", c.ConnID)
	defer fmt.Printf("connID = %d Reader is eixt, remote addr is %s", c.ConnID, c.RemoteAddr())
	defer c.Stop()

	for {
		// 读取客户端的数据到 buf 中，最大 MaxPackageSize 字节
		buf := make([]byte, utils.Config.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("recv buf err: %s", err)
		}

		// 得到当前连接数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

func (c *Connection) Start() {
	fmt.Printf("Conn start...ConnID=%d\n", c.ConnID)
	// 启动当前连接的读数据业务
	c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Printf("Conn stop... connID = %d\n", c.ConnID)
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	err := c.Conn.Close()
	if err != nil {
		return
	}

}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}
