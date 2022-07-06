package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leiTCP/Leinx/liface"
)

/**
存储一切有关此框架的全局参数，供其他模块使用
一些参数是json由用户配置
**/
type Globalobj struct {
	/**
		Server
	**/
	TcpServer liface.IServer //当前全局的Server对象
	Host      string         //当前服务器主机监听的ip
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称
	/**
	leinx
	**/
	Version        string //当前leinx的版本号
	MaxConn        int    //最大连接数
	MaxPackageSize uint32 //当前最大数据包的最大值
	WorkPoolSize   uint32 //当前业务工作的Work池的Goroutine数量
	MaxWorkTaskNum uint32 //每个消息队列能够处理任务的数量
}

/**
	定义一个全局的对外globalobj
**/
var GlobalObject *Globalobj

/*
	加载用户自定义的参数
*/
func (g *Globalobj) Reload() {
	fmt.Println("reload --------------")
	data, err := ioutil.ReadFile("E:\\code\\leiTCP\\Leinx\\conf\\leinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件解析到struct中
	json.Unmarshal(data, &GlobalObject)

}

/**
 提供一个init方法，初始化当前的globalobject
**/
func init() {
	//如果配置文件没有加载
	GlobalObject = &Globalobj{
		Name:           "leiServer app",
		Version:        "tcp4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkPoolSize:   10,
		MaxWorkTaskNum: 1024,
	}
	//GlobalObject.Reload()
}
