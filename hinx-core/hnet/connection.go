package hnet

import (
	"errors"
	"hinx/hinx-core/hconf"
	"hinx/hinx-core/hiface"
	"hinx/hinx-core/hlog"
	"io"
	"net"
	"sync"
	"time"
)

type Connection struct {
	Server    hiface.IServer
	Conn      net.Conn
	IPVersion string
	ConnID    uint32
	isClosed  bool
	ExitChan  chan bool
	msgChan   chan []byte
	MsgHandle hiface.IMsgHandle

	// a set of connection properties
	property     map[string]interface{}
	propertyLock sync.RWMutex

	heartBeater *time.Timer
}

func (c *Connection) Start() {
	hlog.Ins().InfoF("Connection Start... ConnID= %d", c.ConnID)
	hlog.Ins().InfoF("start!", time.Now())
	go c.StartReader()
	go c.StartWriter()

	c.Server.CallOnConnStart(c)

	select {
	case <-c.TickHeartBreaker():
		hlog.Ins().InfoF("timeout!", time.Now())
		c.Stop()
	}
}

func (c *Connection) Stop() {

	hlog.Ins().InfoF("Connection Stop... ConnID= %d", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true

	c.Server.CallOnConnStop(c)
	c.StopHeartBreaker()
	c.Conn.Close()

	c.ExitChan <- true
	c.Server.GetConnManager().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetConnection() net.Conn {
	return c.Conn
}

func (c *Connection) Send(data []byte) error {

	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	_, err := c.Conn.Write(data)
	if err != nil {

		return err
	}

	return nil
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
func NewConnection(parent hiface.IServer, conn net.Conn, connID uint32, ipVersion string, msgHandler hiface.IMsgHandle) *Connection {
	c := &Connection{
		Server:    parent,
		Conn:      conn,
		ConnID:    connID,
		IPVersion: ipVersion,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),

		// no cache channel, used for write & read message between goroutine
		msgChan:   make(chan []byte, 0),
		MsgHandle: msgHandler,

		property: map[string]interface{}{},

		heartBeater: time.NewTimer(time.Duration(hconf.GlobalObject.Heartbeat) * time.Millisecond),
	}

	c.Server.GetConnManager().Add(c)
	return c
}

// StartWriter write message goroutine
func (c *Connection) StartWriter() {

	hlog.Ins().InfoF("writer goroutine is running...")
	//defer fmt.Println("[Writer is exit!], connID=", c.ConnID, "Writer is exit, remote addr is ", c.GetRemoteAddr().String())

	// block
	for {
		select {
		case data := <-c.msgChan:
			{
				c.ResetHeartBreaker()
				_, err := c.Conn.Write(data)
				if err != nil {
					hlog.Ins().ErrorF("send data error %v", err)
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
	hlog.Ins().InfoF("Reader Goroutine is running...")

	defer hlog.Ins().InfoF("[Reader is exit!], connID= %d Reader is exit, remote addr is %s", c.ConnID, c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// create a pack and unpack object

		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())

		//_, err := io.ReadFull(c.GetConnection(), headData)
		_, _, err := c.read(headData)
		if err != nil {
			hlog.Ins().ErrorF("read msg head error: %v", err)
			break
		}

		// reset timer
		c.ResetHeartBreaker()

		msg, err := dp.Unpack(headData)
		if err != nil {
			hlog.Ins().ErrorF("unpack error: %v", err)
			break
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.GetConnection(), data)
			if err != nil {
				hlog.Ins().ErrorF("read msg data error: %v", err)
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

func (c *Connection) read(buf []byte) (int, net.Addr, error) {
	switch c.IPVersion {
	case "tcp4":
		{
			n, err := io.ReadFull(c.GetConnection(), buf)
			//if err != nil {
			//	return 0, nil, err
			//}
			return n, nil, err
		}
	case "udp":
		{
			conn := c.GetConnection().(*net.UDPConn)
			var n int
			var raddr *net.UDPAddr
			var err error
			for n != 0 {
				n, raddr, err = conn.ReadFromUDP(buf[0:])
				if err != nil {
					return 0, nil, err
				}
			}
			return n, raddr, err

		}
	default:
		{
			n, err := io.ReadFull(c.GetConnection(), buf)
			//if err != nil {
			//	return 0, nil, err
			//}
			return n, nil, err
		}
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
		hlog.Ins().ErrorF("Pack error msg id = %d", msgId)
		return errors.New("Pack error msg")
	}

	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) ResetHeartBreaker() {
	c.heartBeater.Reset(time.Duration(hconf.GlobalObject.Heartbeat) * time.Millisecond)
}

func (c *Connection) StopHeartBreaker() {
	c.heartBeater.Stop()
}

func (c *Connection) TickHeartBreaker() <-chan time.Time {
	return c.heartBeater.C
}
