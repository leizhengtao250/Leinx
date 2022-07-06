package lnet

import (
	"fmt"
	"leiTCP/Leinx/liface"
	"leiTCP/Leinx/utils"
	"net"
)

//定义一个server服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip版本
	IpVersion string
	//Ip地址
	Ip string
	//服务器监听的端口
	Port int
	//给当前的server添加一个router，server注册的连接对应的处理业务
	MsgH liface.IMsgHandle
	//该server的连接管理器
	ConnManager liface.IConnManager
	//创建连接之前的hook方法
	OnConnStart func(conn liface.IConnection)
	//销毁连接之前的hook方法
	OnConnStop func(conn liface.IConnection)
}

// //定义客户端连接绑定的handle，应该由用户来自定义
// func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
// 	//回显的业务
// 	fmt.Println("[Conn Handle] callbacktoClient..")
// 	if _, err := conn.Write(data[:cnt]); err != nil {
// 		fmt.Println("write back error", err)
// 		return errors.New("CallBackToClient error")
// 	}

// 	return nil

// }

func (s *Server) Start() {
	fmt.Printf("[start] Server Listenner at IP:%s,Port %d,is starting\n", s.Ip, s.Port)
	//开启工作池
	s.MsgH.StartWorkerPool()
	//1.获取一个tcp的地址
	addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("resolce tcp addr error:", err)
		return
	}
	//2.监听服务器的地址
	listenner, err := net.ListenTCP(s.IpVersion, addr)
	if err != nil {
		fmt.Println("listen", s.IpVersion, "err", err)
		return
	}
	fmt.Println("start leinix", s.Name, "success Listenning")
	var cid uint32 = 0
	//3.阻塞的等待客户端的连接，处理客户端连接业务
	for {
		//如果有客户端连接过来，阻塞会返回
		conn, err := listenner.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err", err)
			continue

		}
		//设置最大连接个数的判断，如果超过最大连接，就关闭新的连接
		if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
			//TODO 给客户端响应一个 错误包
			fmt.Println("[Conn num] too many!!!Please reduce")
			conn.Close()
			continue
		}

		// //客户端已经连接上了，做一些业务，做一个最基本的512字节长度的回显业务
		// go func() {

		// 	for {
		// 		buf := make([]byte, 512)
		// 		cnt, err := conn.Read(buf)
		// 		if err != nil {
		// 			fmt.Println("recv buf err", err)
		// 			continue
		// 		}
		// 		if _, err := conn.Write(buf[:cnt]); err != nil {
		// 			fmt.Println("write back buf err", err)
		// 			continue
		// 		}
		// 	}
		// }()
		cn := NewConnection(s, conn, cid, s.MsgH)
		cid++
		go cn.Start()

	}

}

func (s *Server) Stop() {
	//TODO 将一些服务器的资源，状态或者一些已经开辟的连接信息，进行停止或者回收
	fmt.Println("[STOP] server name:", s.Name)
	s.ConnManager.Clear()
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外服务

	//阻塞状态
}

/**
初始化Server 模块方法
**/
func NewServer() *Server {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IpVersion:   utils.GlobalObject.Version,
		Port:        utils.GlobalObject.TcpPort,
		Ip:          utils.GlobalObject.Host,
		MsgH:        NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}

func (s *Server) GetConnMgr() liface.IConnManager {
	return s.ConnManager
}

//注册OnConnstart
func (s *Server) SetOnConnStart(hookFunc func(conn liface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册OnConnStop
func (s *Server) SetOnConnStop(hookFunc func(conn liface.IConnection)) {
	s.OnConnStart = hookFunc
}

//调用OnConnStart
func (s *Server) CallOnConnStart(conn liface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("------->Call OnConnStart()....")
		s.OnConnStart(conn)
	}

}

//调用OnConnStop
func (s *Server) CallOnConnStop(conn liface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("----->Call OnConnStop()....")
		s.OnConnStop(conn)
	}
}
