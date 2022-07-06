package liface

//定义一个服务器接口
type IServer interface {
	//start server
	Start()
	//stop server
	Stop()
	//do server
	Serve()
	//获取当前的连接管理器
	GetConnMgr() IConnManager

	//注册OnConnstart
	SetOnConnStart(func(conn IConnection))
	//注册OnConnStop
	SetOnConnStop(func(conn IConnection))
	//调用OnConnStart
	CallOnConnStart(conn IConnection)
	//调用OnConnStop
	CallOnConnStop(conn IConnection)
}
