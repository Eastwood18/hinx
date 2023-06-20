package hnet

import (
	"errors"
	"fmt"
	"hinx/hinx-core/hiface"
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
	fmt.Println("connection add ", conn.GetConnID(), "to ConnManager successfully: conn num= ", c.Len())
}

func (c *ConnManager) Remove(conn hiface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.conns, conn.GetConnID())
	fmt.Println("connection delete ", conn.GetConnID(), "to ConnManager successfully: conn num= ", c.Len())
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
	fmt.Println("Clear All connection success", c.Len())
}
