package znet

/*
@author: shg
@since: 2023/1/16 2:46 PM
@mail: shgang97@163.com
*/

type Message struct {
	// 消息的ID
	Id uint32
	// 消息的长度
	DataLen uint32
	// 消息的内容
	Data []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
