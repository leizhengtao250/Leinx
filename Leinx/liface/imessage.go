package liface

/**
 将请求的消息封装到一个message中，定义抽象接口
**/
type IMessage interface {
	//获取消息的ID
	GetMsgID() uint32
	//获取消息的长度
	GetMsgLen() uint32
	//获取消息的内容
	GetData() []byte
	//设置消息ID
	SetMsID(uint32)
	//设置消息长度
	SetMsgLen(uint32)
	//设置消息内容
	SetData([]byte)
}
