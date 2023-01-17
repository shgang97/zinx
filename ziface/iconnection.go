package ziface

import "net"

/*
@author: shg
@since: 2023/1/15
@desc: 定义连接模块的抽象层
*/

type IConnection interface {
	Start()                                            // 启动连接 让当前的连接准备开始工作
	Stop()                                             // 停止连接，结束当前连接的工作
	GetTCPConnection() *net.TCPConn                    // 获取当前连接绑定的 socket conn
	GetConnID() uint32                                 // 获取当前连接的连接id
	RemoteAddr() net.Addr                              // 获取远程客户端的 TCP状态 IP port
	SendMsg(msgId uint32, data []byte) error           // 发送数据，将数据发送给远程的客户端
	SetProperty(propName string, property interface{}) // 设置连接属性
	GetProperty(propName string) (interface{}, error)  // 获取连接属性
	RemoveProperty(propName string)                    // 移除连接属性

}

type HandleFunc func(*net.TCPConn, []byte, int) error // 处理连接业务的方法
