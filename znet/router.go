package znet

import (
	"zinx/ziface"
)

/*
@author: shg
@since: 2023/1/16 12:48 AM
@mail: shgang97@163.com
*/

// BaseRouter 实现 router 时，可以先嵌入这个 BaseRouter 基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

// 这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化

func (router BaseRouter) PreHandle(request ziface.IRequest) {
}

func (router BaseRouter) Handle(request ziface.IRequest) {
}

func (router BaseRouter) PostHandle(request ziface.IRequest) {
}
