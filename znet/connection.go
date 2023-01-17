package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

/*
@author: shg
@since: 2023/1/15
@desc: //TODO
*/

type Connection struct {
	// 连接所在的 server
	TcpServer ziface.IServer
	// 当前连接 socket tcp 套接字
	Conn *net.TCPConn
	// 连接的 ID
	ConnID uint32
	// 当前的连接状态
	IsClosed bool
	// 告知当前连接已经退出/停止的 channel
	ExitChan chan bool
	// 无缓冲管道，用于读、写 Goroutine 之间读消息通信
	MsgChan chan []byte
	// 该连接处理的方法 router
	MsgHandler ziface.IMsgHandler
	// 连接属性集合
	properties map[string]interface{}
	// 保护连接属性集合的锁
	propLock sync.RWMutex
}

func (c *Connection) SetProperty(propName string, property interface{}) {
	c.propLock.Lock()
	defer c.propLock.Unlock()
	c.properties[propName] = property
}

func (c *Connection) GetProperty(propName string) (interface{}, error) {
	c.propLock.RLock()
	defer c.propLock.RUnlock()
	if property, ok := c.properties[propName]; ok {
		return property, nil
	}
	return nil, errors.New(fmt.Sprintf("property[%s] not FOUND", propName))
}

func (c *Connection) RemoveProperty(propName string) {
	c.propLock.Lock()
	defer c.propLock.Unlock()
	delete(c.properties, propName)
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		MsgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		MsgHandler: handler,
		properties: make(map[string]interface{}),
	}
	server.GetConnMgr().Add(c)
	return c
}

// StartReader 读业务方法
func (c *Connection) StartReader() {
	fmt.Printf("Conn[id = %d] start reading\n", c.ConnID)
	defer fmt.Printf("Conn[id = %d] Reader exit\n", c.ConnID)
	defer c.Stop()

	for {
		pack := NewDataPack()
		// 读取客户端的消息头 二进制流 8字节
		headData := make([]byte, pack.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Printf("read msg head error: %s\n", err)
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
		if utils.Config.WorkerPoolSize > 0 {
			// 已经开启工作池，将消息交给工作池进行处理
			c.MsgHandler.SendMsgToRequestChan(&req)
		} else {
			// 执行注册的路由方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// StartWriter 写业务方法
func (c *Connection) StartWriter() {
	fmt.Printf("Conn[id = %d] start writing\n", c.ConnID)
	defer fmt.Printf("Conn[id = %d] writer exit!\n", c.ConnID)

	// 阻塞等待channel 的消息，进行写给客户端
	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Printf("write data error: %s\n", err)
				return
			}
		case <-c.ExitChan:
			// reader 已经退出，此时writer也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Printf("Conn[id = %d] start connecting\n", c.ConnID)
	// 启动当前连接的读数据业务
	go c.StartReader()
	// 启动当前连接的写数据业务
	go c.StartWriter()
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Printf("Conn[id = %d] stop\n", c.ConnID)
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	// 调用连接销毁之前的 hook 方法
	c.TcpServer.CallOnConnStop(c)
	// 关闭 socket 连接
	err := c.Conn.Close()
	if err != nil {
		return
	}
	// 将 conn 从连接管理器中移除
	c.TcpServer.GetConnMgr().Remove(c)
	// 告知 writer 关闭
	c.ExitChan <- true
	// 回收资源
	close(c.ExitChan)
	close(c.MsgChan)

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
	c.MsgChan <- binaryMsg
	return nil
}
