package lnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"leiTCP/Leinx/liface"
	"leiTCP/Leinx/utils"
)

//封包，拆包的模块
type DataPack struct{}

//拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包的长度方法
func (d *DataPack) GetHeadLen() uint32 {
	//datalen uint32 4byte
	//dataid uint32 4byte
	return 8
}

//封包方法
func (d *DataPack) Pack(msg liface.IMessage) ([]byte, error) {
	//创建一个存放byte字节流的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将datalen 写入buf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将dataId 写入buf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	//将message吸入buf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包方法 (将包的head信息读出来，再根据head信息的data长度，再进行一次读)
func (d *DataPack) UnPack(binaryData []byte) (liface.IMessage, error) {
	//创建一个存放byte字节流的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息，得到dataLen 和dataID
	msg := &Message{}

	//读datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Datalen); err != nil {
		return nil, err
	}
	//读dataID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.Datalen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv!")
	}

	return msg, nil
	//读msg内容

}
