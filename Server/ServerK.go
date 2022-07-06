package main

import (
	"fmt"
	"leiTCP/Leinx/liface"
	"leiTCP/Leinx/lnet"
)

/**
  自定义路由
**/

type PingRouter struct {
	lnet.BaseRouter
}

// func (p *PingRouter) PreHandle(request liface.IRequest) {
// 	fmt.Println("call router prehandle")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ping \n"))
// 	if err != nil {
// 		fmt.Println("call back before ping error")
// 	}
// }

func (p *PingRouter) Handle(request liface.IRequest) {
	fmt.Println("call router handle")

	err := request.GetConnection().SendMsg([]byte("hello handle...."), 1)
	if err != nil {
		fmt.Println("call back handle  ping error")
	}
}

// func (p *PingRouter) PostHandle(request liface.IRequest) {
// 	fmt.Println("call router posthandle")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("post ping ping \n"))
// 	if err != nil {
// 		fmt.Println("call back post ping error")
// 	}
// }

/**
服务器端应用程序
**/
func main() {
	/*
		1.创建一个server句柄，使用api
		2.启动server
	*/
	s := lnet.NewServer()
	//添加一个router
	u1 := &lnet.UserRouter1{}
	u2 := &lnet.UserRouter2{}
	s.MsgH.AddRouter(1, u1)
	s.MsgH.AddRouter(2, u2)
	//ping
	s.Serve()
}
