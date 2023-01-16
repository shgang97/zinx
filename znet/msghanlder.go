package znet

import (
	"fmt"
	"strconv"
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
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{Apis: make(map[uint32]ziface.IRouter)}
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
