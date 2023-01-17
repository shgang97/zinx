package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

/*
@author: shg
@since: 2023/1/17 3:28 PM
@mail: shgang97@163.com
*/

/*
ConnManager
连接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                  // 保护连接结合的读写锁
}

func (manager *ConnManager) Add(conn ziface.IConnection) {
	// 添加连接加写锁
	manager.connLock.Lock()
	//defer manager.connLock.Unlock()
	// 将 conn 加入到 connections
	manager.connections[conn.GetConnID()] = conn
	// 打印日志不需要枷锁保护
	manager.connLock.Unlock()
	fmt.Printf("add conn[id=%d] into connections\n", conn.GetConnID())
}

func (manager *ConnManager) Remove(conn ziface.IConnection) {
	// 删除连接加写锁
	manager.connLock.Lock()
	delete(manager.connections, conn.GetConnID())
	manager.connLock.Unlock()
	fmt.Printf("delete conn[id=%d] from connections\n", conn.GetConnID())
}

func (manager *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	// 获取连接加读锁
	manager.connLock.RLock()
	defer manager.connLock.RUnlock()
	if conn, ok := manager.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New(fmt.Sprintf("conn[id=%d] is not FOUND", connId))
	}
}

func (manager *ConnManager) Size() int {
	return len(manager.connections)
}

func (manager *ConnManager) Clear() {
	manager.connLock.Lock()
	defer manager.connLock.Unlock()
	// 删除并停止conn工作
	for connId, conn := range manager.connections {
		conn.Stop()
		// 清空 connections 可以将所有 connection 停止，然后将 connections 指向一个新的 map
		delete(manager.connections, connId)
	}
}

func NewConnManager() *ConnManager {
	return &ConnManager{connections: make(map[uint32]ziface.IConnection)}
}
