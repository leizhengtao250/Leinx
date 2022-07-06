package liface

/*
	IRequest接口：
	实际上把客户端请求的连接信息和请求的数据包装到了一个Request
*/

type IRequest interface {
	//得到当前连接
	GetConnection() IConnection
	//得到请求的数据
	GetData() []byte
	//得到请求的数据ID
	GetMsgID() uint32
}
