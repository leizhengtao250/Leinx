package lnet

import (
	"errors"
	"fmt"
	"io"
	"leiTCP/Leinx/liface"
	"leiTCP/Leinx/utils"
	"net"
)

/**
	连接模块
**/
type Connection struct {

	//当前连接的socket 套接字
	Conn *net.TCPConn
	//当前连接的ID
	ConnID uint32
	//当前连接的状态
	IsClosed bool
	//当前连接的绑定业务的API
	//handleAPI liface.HandleFunc
	//告知当前连接已经退出的/停止 channel reader告知writer 退出
	ExitChan chan bool
	//该连接处理的方法Router
	Router liface.IMsgHandle
	//无缓冲的通道，用于读写Goroutine之间的消息通信
	msgChan chan []byte
	//当前Conn 隶属那个Server
	TcpServer liface.IServer
}

//初始化连接模块的方法
func NewConnection(CurrServer liface.IServer, conn *net.TCPConn, connId uint32, router liface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: CurrServer,
		Conn:      conn,
		ConnID:    connId,
		//handleAPI: callback_api,
		IsClosed: false,              //当前开启，因此为false
		ExitChan: make(chan bool, 1), //
		Router:   router,
		msgChan:  make(chan []byte),
	}
	//将conn 加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

//从连接中读数据
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...cid=")
	defer fmt.Println("connId=", c.ConnID, "Reader is exit,remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//得到当前连接conn数据的Request请求数据
		//创建一个拆包解包的对象
		dp := NewDataPack()
		//读取客户端的Msg Head二进制流8个字符
		buf := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), buf); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包，得到msgID和msgDataLen 放在msg消息中
		msgHead, err := dp.UnPack(buf)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}
		datalen := msgHead.GetMsgLen()
		msgTemp := make([]byte, datalen)
		if datalen > 0 {
			_, err := io.ReadFull(c.Conn, msgTemp)
			if err != nil {
				fmt.Println("unpack error:", err)
			}
		}
		fmt.Println(string(msgTemp))
		msg := &Message{
			Id:      msgHead.GetMsgID(),
			Datalen: datalen,
			Data:    msgTemp,
		}

		// //调用当前连接所绑定的HandleApi
		// if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		// 	fmt.Println("ConnId", c.ConnID, "handle is error", err)
		// 	break
		// }

		//最终得到当前的conn数据放在request中
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkPoolSize > 0 {
			//将消息交给TaskQueue，由worker来处理

			c.Router.SendMsgToTaskQueue(&req)
		} else {
			//从路由中，找到注册绑定的conn对应的router调用
			go c.Router.DoMsgHandle(&req)
		}

		// data := []byte{'1', '2', 'x', 'x', 'g'}
		// c.SendMsg(data, 4)
		c.SendMsg(msgTemp, msg.GetMsgID())

	}
}

/*
	写消息的goroutine,专门发送给客户端消息的模块
*/

func (c *Connection) startWriter() {
	fmt.Println("[Writer Gortine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	//不断的阻塞的等待channel的消息,进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error:", err)
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也要退出
			return
		}

	}

}

//提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(data []byte, msgId uint32) error {
	if c.IsClosed == true {
		return errors.New("Connection close when send msg")
	}

	dp := NewDataPack()
	//先封装成msg
	dataLen := len(data)
	msg := &Message{}
	msg.SetData(data)
	msg.SetMsgLen(uint32(dataLen))
	msg.SetMsID(msgId)
	sendBuf, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("pack error:", err)
		return err
	}
	// //将数据发送给客户端
	// if _, err := c.Conn.Write(sendBuf); err != nil {
	// 	fmt.Println("Write msg id:", msgId, "error:", err)
	// 	return errors.New("conn send error")
	// }
	//将数据发送给channle
	c.msgChan <- sendBuf
	return nil
}

//启动连接 当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn start()... ConnID=", c.ConnID)
	//从当前连接中的读数据业务
	go c.StartReader()
	//从当前连接中的写数据业务
	go c.startWriter()

	//按照开发者传递进来的创建连接之后要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

//停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn stop()... ConnID=", c.ConnID)
	// 如果当前连接已经关闭
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	//在销毁连接之前需要执行业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket连接
	c.Conn.Close()
	//将当前conn从ConnMgr中除掉
	c.TcpServer.GetConnMgr().Remove(c)
	//告知writer关闭
	c.ExitChan <- true
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

//获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的连接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP 状态 IP PORT
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
