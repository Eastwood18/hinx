package hnet

import (
	"errors"
	"hinx/hinx-core/hiface"
	"hinx/hinx-core/hlog"
	"sync"
)

type ConnManager struct {
	conns    map[uint32]hiface.IConnection
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		conns: map[uint32]hiface.IConnection{},
	}
}

func (c *ConnManager) Add(conn hiface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.conns[conn.GetConnID()] = conn
	hlog.Ins().InfoF("connection add %d to ConnManager successfully: conn num= %d ", conn.GetConnID(), c.Len())
}

func (c *ConnManager) Remove(conn hiface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.conns, conn.GetConnID())
	hlog.Ins().InfoF("connection delete %d to ConnManager successfully: conn num= %d ", conn.GetConnID(), c.Len())
}

func (c *ConnManager) Get(connID uint32) (hiface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	connection, ok := c.conns[connID]
	if !ok {
		return nil, errors.New("connection not found")
	}
	return connection, nil
}

func (c *ConnManager) Len() int {
	return len(c.conns)
}

func (c *ConnManager) ClearConns() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for _, conn := range c.conns {
		conn.Stop()
	}
	c.conns = make(map[uint32]hiface.IConnection, 0)
	hlog.Ins().InfoF("Clear All connection success %d", c.Len())
}
