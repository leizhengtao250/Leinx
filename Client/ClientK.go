package main

import (
	"fmt"
	"io"
	"leiTCP/Leinx/liface"
	"leiTCP/Leinx/lnet"
	"net"
	"time"
)

//模拟客户端
func main() {
	fmt.Println("client start")
	//1.直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error", err)
		return
	}

	for {
		//2.发送封包的message消息
		dp := lnet.NewDataPack()
		msg1 := lnet.NewMessage(1, []byte("hello I am client3,msgid=1"))
		msg2 := lnet.NewMessage(2, []byte("hello I am client4,msgid=2"))
		// sendbuf, err := dp.Pack(msg)
		// if err != nil {
		// 	fmt.Println("pack error:", err)
		// 	continue
		// }
		// if _, err := conn.Write(sendbuf); err != nil {
		// 	fmt.Println("write error:", err)
		// }

		// //cpu 阻塞
		// time.Sleep(1 * time.Second)
		// recieveHead := make([]byte, dp.GetHeadLen())
		// if _, err := io.ReadFull(conn, recieveHead); err != nil {
		// 	fmt.Println("read error:", err)
		// 	return
		// }
		// msgU, err := dp.UnPack(recieveHead)
		// if err != nil {
		// 	fmt.Println("unpack error:", err)
		// 	return
		// }
		// datalen := msgU.GetMsgLen()
		// msgData := make([]byte, datalen)
		// if _, err := io.ReadFull(conn, msgData); err != nil {
		// 	fmt.Println("")
		// 	return
		// }
		// fmt.Println(string(msgData))
		Send(conn, msg1, dp)
		Send(conn, msg2, dp)
		//cpu 阻塞
		time.Sleep(1 * time.Second)
		recieveHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, recieveHead); err != nil {
			fmt.Println("read error:", err)
			return
		}
		msgU, err := dp.UnPack(recieveHead)
		if err != nil {
			fmt.Println("unpack error:", err)
			return
		}
		datalen := msgU.GetMsgLen()
		msgData := make([]byte, datalen)
		if _, err := io.ReadFull(conn, msgData); err != nil {
			fmt.Println("")
			return
		}
		fmt.Println(string(msgData))
	}

}

func Send(conn net.Conn, msg liface.IMessage, dp liface.IDataPack) {
	sendbuf, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("pack error:", err)
	}

	if _, err := conn.Write(sendbuf); err != nil {
		fmt.Println("write error:", err)
	}
}
