package lnet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

//只是负责datapack 拆包封包单元测试
func TestDataPack(t *testing.T) {
	fmt.Println("----------------------------------------------")
	/**
		模拟的服务器
	**/
	//1.创建socketTCP
	listen, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err", err)
	}
	//创建一个go承载从客户端读取业务
	go func() {
		//2.从客户端读取数据，拆包处理
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("server accpet error:", err)
			}
			go func(con net.Conn) {
				dp := NewDataPack()
				for {
					//处理客户端的请求
					//------>拆包的过程
					headData := make([]byte, dp.GetHeadLen())
					//先把datalen读出来
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error:", err)
					}
					dataLen := msgHead.GetMsgLen()
					msg := make([]byte, dataLen)
					if dataLen > 0 {
						//msg有数据
						io.ReadFull(conn, msg)

					}
					fmt.Println(string(msg))
					//dataId := msg.GetMsgID()
					//再把消息内容读出来

				}
			}(conn)

		}
	}()

	//模拟客户端
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:7777")
		if err != nil {
			fmt.Println("conn client error:", err)
			return
		}
		//创建一个封包对象
		dp := NewDataPack()
		//模拟粘包问题，封装两个msg 一起发送
		//封装第一个msg1包
		msg1 := &Message{
			Id:      1,
			Datalen: 5,
			Data:    []byte{'h', 'e', 'l', 'l', 'o'},
		}
		sendData1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("client pack msg1 error:", err)
			return
		}
		//封装第二个msg2包
		msg2 := &Message{
			Id:      2,
			Datalen: 6,
			Data:    []byte{'o', 'r', 'a', 'n', 'g', 'e'},
		}
		sendData2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("client pack msg1 error:", err)
			return
		}
		//将两个包黏在一起
		sendData1 = append(sendData1, sendData2...)
		//一次性发给客户端
		conn.Write(sendData1)
	}()
	//客户端阻塞
	select {
	case <-time.After(time.Second):
		return
	}

}
