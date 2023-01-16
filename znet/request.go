package znet

import (
	"zinx/ziface"
)

/*
@author: shg
@since: 2023/1/16
*/

type Request struct {
	// 与客户端建立的连接
	conn ziface.IConnection
	// 客户端发送的请求数据
	msg ziface.IMessage
}

func (request Request) GetConnection() ziface.IConnection {
	return request.conn
}

func (request Request) GetData() []byte {
	return request.msg.GetData()
}

func (request Request) GetDataLen() uint32 {
	return request.msg.GetDataLen()
}

func (request Request) GetMsgId() uint32 {
	return request.msg.GetMsgId()
}
