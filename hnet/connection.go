package hnet

import (
	"errors"
	"fmt"
	"hinx/hiface"
	"io"
	"net"
	"sync"
)

type Connection struct {
	TcpServer hiface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	ExitChan  chan bool
	msgChan   chan []byte
	MsgHandle hiface.IMsgHandle

	// a set of connection properties
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value

}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	val, ok := c.property[key]
	if !ok {
		return nil, errors.New("Property not found")
	}
	return val, nil
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// NewConnection init Connection module function
func NewConnection(parent hiface.IServer, conn *net.TCPConn, connID uint32, msgHandler hiface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: parent,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),

		// no cache channel, used for write & read message between goroutine
		msgChan:   make(chan []byte, 0),
		MsgHandle: msgHandler,

		property: map[string]interface{}{},
	}

	c.TcpServer.GetConnManager().Add(c)
	return c
}

// StartWriter write message goroutine
func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine is running...")
	defer fmt.Println("[Writer is exit!], connID=", c.ConnID, "Writer is exit, remote addr is ", c.GetRemoteAddr().String())

	// block
	for {
		select {
		case data := <-c.msgChan:
			{
				_, err := c.Conn.Write(data)
				if err != nil {
					fmt.Println("send data error", err)
					continue
				}
			}
		case <-c.ExitChan:
			{
				return
			}
		}
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("[Reader is exit!], connID=", c.ConnID, "Reader is exit, remote addr is ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// create a pack and unpack object
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error: ", err)
				break
			}
		}
		msg.SetMsgData(data)

		req := &Request{
			conn: c,
			msg:  msg,
		}

		c.MsgHandle.SendMsgToTaskQueue(req)
	}
}
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed")
	}

	// pack data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	c.msgChan <- binaryMsg
	return nil
}
func (c *Connection) Start() {
	fmt.Println("Connection Start... ConnID=", c.ConnID)
	go c.StartReader()
	go c.StartWriter()

	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop... ConnID= ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true

	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()

	c.ExitChan <- true
	c.TcpServer.GetConnManager().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
