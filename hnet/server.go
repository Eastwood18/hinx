package hnet

import (
	"fmt"
	"hinx/hiface"
	"hinx/utils"
	"net"
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
	fmt.Println("Add router success")
}

func (s *Server) Start() {
	/*
		1. get a TCP Addr
		2. listen server addr
		3. blocking wait client connection, handle client business(i/o)
	*/
	fmt.Printf("[Hinx] Server Name: %s, listener at ip: %s is starting\n",
		utils.GlobalObject.Name, fmt.Sprintf("%s:%d", utils.GlobalObject.Host, utils.GlobalObject.TcpPort))

	go func() {
		s.MsgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp addr error: ", err)
			return
		}
		fmt.Println("start Hinx server success", s.Name, "Listening...")
		var cid uint32

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//set Maximal connections judge
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too many connections")
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Hinx server name", s.Name)

	s.ConnManager.ClearConns()
}
func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) GetConnManager() hiface.IConnManager {
	return s.ConnManager
}

// NewServer init server
func NewServer(name string) hiface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}

	return s
}
