package snet

import (
	"fmt"
	"game_server_silk/siface"
	"github.com/pkg/errors"
	"sync"
)

/*
	连接管理模块
 */
type ConnManager struct {
	//管理的连接集合
	connections map[uint32]siface.IConnection  
	
	//保护连接集合的读写锁
	connLock sync.RWMutex
} 

//初始化连接管理集合的方法
func NewConnManager() *ConnManager  {
	return &ConnManager{
		connections: make(map[uint32] siface.IConnection),
	}
}

//添加连接
func (connMgr *ConnManager) Add(conn siface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn 加入到ConnManager连接管理集合中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("[客户端连接成功加入到连接管理集合中]connID = ", conn.GetConnID(),
		" add to ConnManager successfully: conn num = ", connMgr.Len())
}

//删除连接
func (connMgr *ConnManager) Remove(conn siface.IConnection){
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接管理集合中的连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("[从连接管理集合中成功删除客户端连接]connID = ", conn.GetConnID(),
		" remove from ConnManager successfully: conn num = ", connMgr.Len())
}
//根据connID获取连接
func (connMgr *ConnManager) Get(connID uint32) (siface.IConnection, error) {
	//保护共享资源map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	//从连接管理集合中，找到连接
	if conn, ok := connMgr.connections[connID]; ok {
		//找到了
		return conn, nil
	} else {
		return nil, errors.New("[没有找到该客户端连接]connection not found!")
	}
}
//得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}
//清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn工作
	for connID, conn := range connMgr.connections {
		//停止客户端连接
		conn.Stop()
		//从连接管理集合中，删除所有客户端连接
		delete(connMgr.connections, connID)
	}
	fmt.Println("[删除连接管理集合中所有连接]Clear All connections succ! conn num = ", connMgr.Len())
}