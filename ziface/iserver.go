package ziface

/*
@author: shg
@since: 2023/1/15
@desc: 提供Server抽象层全部接口声明
*/

type IServer interface {
	Start()                                 //启动服务器方法
	Stop()                                  //停止服务器方法
	Serve()                                 //开启业务服务方法
	AddRouter(msgId uint32, router IRouter) // 路由功能给当前的服务注册一个路由方法，共客户端连接处理使用
}
