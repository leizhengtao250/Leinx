package liface

/**
路由抽象接口
路由里的数据都是IRqueset
**/
type IRouter interface {
	//在处理conn连接之前的方法
	PreHandle(request IRequest)
	//在处理conn连接之后的方法
	PostHandle(request IRequest)
	//在处理conn连接的主方法
	Handle(request IRequest)
}
