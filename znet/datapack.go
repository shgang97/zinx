package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

/*
@author: shg
@since: 2023/1/16 3:00 PM
@mail: shgang97@163.com
*/

/*
DataPack
封包、拆包 具体实现
*/
type DataPack struct {
}

func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

func (pack *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 + Id uint32
	return 8
}

/*
Pack
封包
*/
func (pack *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	// 将 dataLen 写进 dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// 将 MsgId 写进 dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将 data 数据写进到 dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

/*
Unpack
拆包
*/
func (pack *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(data)
	// 只解压 head 信息，得到 dataLen 和 MsgId
	msg := &Message{}
	// 读 dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读 MsgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	if 0 < utils.Config.MaxPackageSize && utils.Config.MaxPackageSize < msg.GetDataLen() {
		return nil, errors.New("too large data recv")
	}
	return msg, nil
}
