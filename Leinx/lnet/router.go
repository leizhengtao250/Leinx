package lnet

import "leiTCP/Leinx/liface"

//实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写就可以了
type BaseRouter struct{}

//这里的三个方法默认为空，按照需要继承实现
//在处理conn连接之前的方法
func (br *BaseRouter) PreHandle(request liface.IRequest) {

}

//在处理conn连接之后的方法
func (br *BaseRouter) PostHandle(request liface.IRequest) {

}

//在处理conn连接的主方法
func (br *BaseRouter) Handle(request liface.IRequest) {

}
