package lnet

import "leiTCP/Leinx/liface"

type Request struct {
	//已经和客户端建立好的连接
	conn liface.IConnection

	//客户端请求的数据
	msg liface.IMessage
}

func (r *Request) GetConnection() liface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
