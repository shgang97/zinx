package ziface

/*
@author: shg
@since: 2023/1/16 6:15 PM
@mail: shgang97@163.com
*/

/*
IMsgHandler
消息管理抽象层
*/
type IMsgHandler interface {
	DoMsgHandler(request IRequest)          // 调度/执行对应的 router 处理小
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
	StartWorkerPool()                       // 启动 worker 工作池
	SendMsgToRequestChan(request IRequest)  //将请求发送给任务请求channel进行处理
}
