package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
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
	defer fmt.Printf("connID = %d Reader is exit, remote addr is %s", c.ConnID, c.RemoteAddr())
	defer c.Stop()

	for {
		pack := NewDataPack()
		// 读取客户端的消息头 二进制流 8字节
		headData := make([]byte, pack.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Printf("read msg head error: %s", err)
			break
		}
		// 拆包
		msg, err := pack.Unpack(headData)
		if err != nil {
			fmt.Printf("unpack head error: %s", err)
			break
		}
		// 根据 dataLen 再次读取 Data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Printf("read msg data error: %s", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前连接数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}
	// 将 data 进行封包
	pack := NewDataPack()
	binaryMsg, err := pack.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Printf("Pack error msg id = %d", msgId)
		return errors.New("pack error msg")
	}
	// 将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Printf("Write msg id[%d] error: %s", msgId, err)
		return errors.New("conn write error")
	}
	return nil
}
