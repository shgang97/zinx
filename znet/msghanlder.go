package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
@author: shg
@since: 2023/1/16 6:18 PM
@mail: shgang97@163.com
*/

/*
MsgHandler
消息处理模块的实现
*/
type MsgHandler struct {
	// msgId - router 映射
	Apis map[uint32]ziface.IRouter
	// 负责 worker 取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作 worker 池的 worker 数量
	WorkerPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.Config.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.Config.WorkerPoolSize),
	}
}

func (handler *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 从Request中找到msgId
	router, ok := handler.Apis[request.GetMsgId()]
	if !ok {
		fmt.Printf("api msgId[%d] is not found! Please Register\n", request.GetMsgId())
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (handler *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断当前 msg 绑定的API处理方法是否已经存在
	if _, ok := handler.Apis[msgId]; ok {
		panic("repeat api, msgId = " + strconv.Itoa(int(msgId)))
	}
	handler.Apis[msgId] = router
	fmt.Printf("Add api msgId[%d] success\n", msgId)
}

/*
StartWorkerPool
启动一个worker池
*/
func (handler *MsgHandler) StartWorkerPool() {
	// 根据 workerPoolSize 分别开启Worker，每个Worker用一个go
	for i := 0; i < int(handler.WorkerPoolSize); i++ {
		// 一个 worker 被启动
		// 1 当前的 worker 对应的 channel 消息队列
		handler.TaskQueue[i] = make(chan ziface.IRequest, utils.Config.MaxWorkerTaskLen)
		// 启动当前的 worker，阻塞等待消息从 channel 中传递进来
		go handler.StartOneWorker(i, handler.TaskQueue[i])
	}
}

/*
StartOneWorker
启动一个worker
*/
func (handler *MsgHandler) StartOneWorker(id int, requestChan chan ziface.IRequest) {
	fmt.Printf("worker[id = %d] is started\n", id)
	// 阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息传递过来，出列的就是一个客户端的request
		case request := <-requestChan:
			handler.DoMsgHandler(request)
		}
	}
}

/*
SendMsgToRequestChan
将请求发送给 requestChan，由worker进行处理
*/
func (handler *MsgHandler) SendMsgToRequestChan(request ziface.IRequest) {
	// 1. 将消息平均分配给不通过的 worker
	// 根据客户端建立的 ConnID 进行分配
	workerId := request.GetConnection().GetConnID() % handler.WorkerPoolSize
	fmt.Printf("Add request[%d] in conn[%d] into requestChan[%d]\n", request.GetMsgId(), request.GetConnection().GetConnID(), workerId)
	// 将 request 发送给 requestChan
	handler.TaskQueue[workerId] <- request
}
