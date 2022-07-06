package lnet

import (
	"errors"
	"fmt"
	"leiTCP/Leinx/liface"
	"sync"
)

/**
连接管理模块
**/

type ConnManager struct {
	connections map[uint32]liface.IConnection //管理的连接集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

//创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]liface.IConnection),
	}
}

//添加连接
func (connMgr *ConnManager) Add(conn liface.IConnection) {
	//保护共享资源 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager sucessfully:conn num=", connMgr.Len())
}

//删除连接
func (connMgr *ConnManager) Remove(conn liface.IConnection) {
	//保护共享资源 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connection:", conn.GetConnID(), "delete successfully")
}

//根据connID获取连接
func (connMgr *ConnManager) Get(connID uint32) (liface.IConnection, error) {
	//保护共享资源 加写锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND")
	}
}

//得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//清除并终止所有连接
func (connMgr *ConnManager) Clear() {
	//保护共享资源 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//删除并停止conn工作
	for connID, conn := range connMgr.connections {
		//conn 停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All connections succ! conn num=", connMgr.Len())
}
