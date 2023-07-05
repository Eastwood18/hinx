package hnet

import (
	"fmt"
	"hinx/hinx-core/hconf"
	"hinx/hinx-core/hiface"
	"hinx/hinx-core/hlog"
	"net"
	"time"
)

// Server implement of IServer, define a Server module
type Server struct {
	Name        string // server name
	IPVersion   string // server ip version
	IP          string // server ip
	Port        int    // server bind port
	MsgHandler  hiface.IMsgHandle
	ConnManager hiface.IConnManager

	OnConnStart func(conn hiface.IConnection)
	OnConnStop  func(conn hiface.IConnection)
}

func (s *Server) SetOnConnStart(f func(connection hiface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(connection hiface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(connection hiface.IConnection) {
	if s.OnConnStart == nil {
		return
	}
	s.OnConnStart(connection)
}

func (s *Server) CallOnConnStop(connection hiface.IConnection) {
	if s.OnConnStart == nil {
		return
	}
	s.OnConnStop(connection)
}

func (s *Server) AddRouter(id uint32, router hiface.IRouter) {
	s.MsgHandler.AddRouter(id, router)
	hlog.Ins().InfoF("Add router success")
}

func (s *Server) Start() {
	/*
		1. get a TCP Addr
		2. listen server addr
		3. blocking wait client connection, handle client business(i/o)
	*/

	hlog.Ins().InfoF("[Hinx] Server Name: %s, listener at ip: %s is starting\n",
		hconf.GlobalObject.Name, fmt.Sprintf("%s:%d", hconf.GlobalObject.Host, hconf.GlobalObject.TcpPort))
	switch s.IPVersion {
	case "tcp4":
		goto TCP
	case "udp":
		goto UDP
	}

TCP:
	go func() {
		s.MsgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			hlog.Ins().ErrorF("resolve tcp addr error: %v", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			hlog.Ins().ErrorF("listen tcp addr error:  %v", err)
			return
		}

		hlog.Ins().InfoF("start Hinx server success %s Listening...", s.Name)
		var cid uint32

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				hlog.Ins().ErrorF("Accept err: %v", err)
				continue
			}
			//set Maximal connections judge
			if s.ConnManager.Len() >= hconf.GlobalObject.MaxConn {
				hlog.Ins().InfoF("Too many connections")
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.IPVersion, s.MsgHandler)
			cid++

			dealConn.Start()
		}
	}()
	select {}
UDP:
	go func() {
		addr, err := net.ResolveUDPAddr("udp", ":8899")
		if err != nil {
			hlog.Ins().ErrorF("resolve udp addr error: %v", err)
			return
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			hlog.Ins().ErrorF("listen udp addr error: %v", err)
			return
		}
		defer conn.Close()

		hlog.Ins().InfoF("udp established!")

		//for {
		//
		//	buf := [1024]byte{}
		//	n, read, err := conn.ReadFromUDP(buf[0:])
		//	if err != nil {
		//		return
		//	}
		//	fmt.Println(fmt.Sprintf("msg from: %s, is: %s", read, string(buf[:n])))
		//
		//	conn.WriteToUDP([]byte("Hello world!"), read)
		//
		//}

		dealConn := NewConnection(s, conn, 0, s.IPVersion, s.MsgHandler)

		dealConn.Start()

	}()

}

func (s *Server) Stop() {
	hlog.Ins().InfoF("[STOP] Hinx server name %s", s.Name)
	s.ConnManager.ClearConns()
}
func (s *Server) Serve() {
	s.Start()

	//select {}
	time.Sleep(10 * time.Second)
}

func (s *Server) GetConnManager() hiface.IConnManager {
	return s.ConnManager
}

// NewServer init server
func NewServer(name string) hiface.IServer {

	s := &Server{
		Name:      hconf.GlobalObject.Name,
		IPVersion: "tcp4",
		//IPVersion:   "udp",
		IP:          hconf.GlobalObject.Host,
		Port:        hconf.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}

	return s
}
