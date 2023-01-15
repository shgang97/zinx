package ziface

/*
@author: shg
@since: 2023/1/16 12:47 AM
@mail: shgang97@163.com
*/

type IRouter interface {
	PreHandle(request IRequest)  // 处理 conn 业务之前的钩子方法 Hook
	Handle(request IRequest)     // 处理 conn 业务的主方法 Hook
	PostHandle(request IRequest) // 处理 conn 业务之后的钩子方法 Hook

}
