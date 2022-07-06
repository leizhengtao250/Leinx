package lnet

import (
	"fmt"
	"leiTCP/Leinx/liface"
)

type UserRouter2 struct {
	BaseRouter
}

func (u *UserRouter2) PreHandle(request liface.IRequest) {
	fmt.Println("userrouter2 prehandle")
}

func (u *UserRouter2) Handle(request liface.IRequest) {
	fmt.Println("userrouter2 handle")
}

func (u *UserRouter2) PostHandle(request liface.IRequest) {
	fmt.Println("userrouter2 posthandle")
}
