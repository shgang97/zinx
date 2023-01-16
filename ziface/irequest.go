package ziface

/*
@author: shg
@since: 2023/1/16
*/

/*
IRequest 接口
把客户端请求的连接信息和请求的数据 封装到 Request 中
*/
type IRequest interface {
	GetConnection() IConnection // 得到当前的请求连接
	GetData() []byte            // 得到当前的请求数据
	GetDataLen() uint32         // 得到当前请求的数据长度
	GetMsgId() uint32           // 得到当前请求的数据 id
}
