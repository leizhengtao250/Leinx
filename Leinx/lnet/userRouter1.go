package lnet

import (
	"fmt"
	"leiTCP/Leinx/liface"
)

type UserRouter1 struct {
	BaseRouter
}

func (u *UserRouter1) PreHandle(request liface.IRequest) {
	fmt.Println("userrouter1 prehandle")
}

func (u *UserRouter1) Handle(request liface.IRequest) {
	fmt.Println("userrouter1 handle")
}

func (u *UserRouter1) PostHandle(request liface.IRequest) {
	fmt.Println("userrouter1 posthandle")
}
