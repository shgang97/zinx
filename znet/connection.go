package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

/*
@author: shg
@since: 2023/1/15
@desc: //TODO
*/

type Connection struct {
	Conn      *net.TCPConn      // 当前连接 socket tcp 套接字
	ConnID    uint32            // 连接的 ID
	IsClosed  bool              // 当前的连接状态
	HandleAPI ziface.HandleFunc // 当前连接所绑定的业务处理方法 API
	ExitChan  chan bool         // 告知当前连接已经退出/停止的 channel
}

func NewConnection(conn *net.TCPConn, connID uint32, handleFunc ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		IsClosed:  false,
		HandleAPI: handleFunc,
		ExitChan:  nil,
	}
	return c
}

// StartReader 读业务方法
func (c *Connection) StartReader() {
	fmt.Printf("Conn start reading..., connID = %d\n", c.ConnID)
	defer fmt.Printf("connID = %d Reader is eixt, remote addr is %s", c.ConnID, c.RemoteAddr())
	defer c.Stop()

	for {
		// 读取客户端的数据到 buf 中，最大 512 字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("recv buf err: %s", err)
		}

		// 调用当前连接所绑定的 HandleAPI
		if err := c.HandleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Printf("ConnID: %d, handle is err: %s\n", c.ConnID, err)
			break
		}
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
	//TODO implement me
	panic("implement me")
}

func (c *Connection) GetConnID() uint32 {
	//TODO implement me
	panic("implement me")
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}
