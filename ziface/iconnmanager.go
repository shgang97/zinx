package ziface

/*
@author: shg
@since: 2023/1/17 3:28 PM
@mail: shgang97@163.com
*/

/*
IConnManager
连接管理模块抽象层
*/
type IConnManager interface {
	Add(conn IConnection)                   // 添加连接
	Remove(conn IConnection)                // 移除连接
	Get(connId uint32) (IConnection, error) // 获取连接
	Size() int                              // 连接总数
	Clear()                                 // 清除所有连接
}
