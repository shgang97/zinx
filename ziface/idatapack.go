package ziface

/*
@author: shg
@since: 2023/1/16 3:00 PM
@mail: shgang97@163.com
*/

/*
IDataPack
封装、拆包 模块
直接面向TCP连接中的数据流，用于处理TCP粘包问题
*/
type IDataPack interface {
	GetHeadLen() uint32                // 获取包的头长度
	Pack(msg IMessage) ([]byte, error) // 封包
	Unpack([]byte) (IMessage, error)   // 拆包
}
