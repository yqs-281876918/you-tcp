package you_net

import (
	"sync"
)

type ConnManager struct {
	connections map[uint32]*YouConnection
	rwMutex     sync.RWMutex
}

func (this *ConnManager) AddConn(conn *YouConnection) {
	this.rwMutex.Lock()
	defer this.rwMutex.Unlock()
	if this.ConnCount() >= 1000 {
		conn.Close()
		return
	}
	conn.Start()
	this.connections[conn.GetConnectionID()] = conn
}

func (this *ConnManager) RemoveConn(conn *YouConnection) {
	this.rwMutex.Lock()
	defer this.rwMutex.Unlock()
	delete(this.connections, conn.GetConnectionID())
}

func (this *ConnManager) GetConn(id uint32) *YouConnection {
	this.rwMutex.RLock()
	defer this.rwMutex.RUnlock()
	if conn, ok := this.connections[id]; ok {
		return conn
	}
	return nil
}

func (this *ConnManager) ConnCount() int {
	return len(this.connections)
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]*YouConnection),
	}
}
