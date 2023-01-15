package ziface

import "net"

/*
@author: shg
@since: 2023/1/15
@desc: 定义连接模块的抽象层
*/

type IConnection interface {
	Start()                         // 启动连接 让当前的连接准备开始工作
	Stop()                          // 停止连接，结束当前连接的工作
	GetTCPConnection() *net.TCPConn // 获取当前连接绑定的 socket conn
	GetConnID() uint32              // 获取当前连接的连接id
	RemoteAddr() net.Addr           // 获取远程客户端的 TCP状态 IP port
	Send(data []byte) error         // 发送数据，将数据发送给远程的客户端
}

type HandleFunc func(*net.TCPConn, []byte, int) error // 处理连接业务的方法
