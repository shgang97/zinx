package ziface

/*
@author: shg
@since: 2023/1/16 2:45 PM
@mail: shgang97@163.com
*/

// IMessage 将请求的消息封装到一个Message中，定义抽象的接口
type IMessage interface {
	GetMsgId() uint32   // 获取消息的Id
	GetDataLen() uint32 // 获取消息的长度
	GetData() []byte    // 获取消息的内容

	SetMsgId(uint32)   // 设置消息的Id
	SetDataLen(uint32) // 设置消息的长度
	SetData([]byte)    // 设置消息的内容
}
