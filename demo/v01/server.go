package main

import (
	"fmt"
	"hinx/hiface"
	"hinx/hnet"
)

type PingRouter struct {
	hnet.BaseRouter
}

//func (p *PingRouter) PreHandle(request hiface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
//	if err != nil {
//		fmt.Println("call back before ping err: ", err)
//
//	}
//
//}

func (p *PingRouter) Handle(request hiface.IRequest) {
	fmt.Println("Call Router Handle...")

	fmt.Println("recv from client: msg ID= ", request.GetMsgID(), ", data= ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping..."))
	if err != nil {
		fmt.Println(err)
		return
	}
}

type HelloRouter struct {
	hnet.BaseRouter
}

func (h *HelloRouter) Handle(request hiface.IRequest) {
	fmt.Println("Call Router Handle...")

	fmt.Println("recv from client: msg ID= ", request.GetMsgID(), ", data= ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello..."))
	if err != nil {
		fmt.Println(err)
		return
	}
}

//func (p *PingRouter) PostHandle(request hiface.IRequest) {
//	fmt.Println("Call Router PostHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
//	if err != nil {
//		fmt.Println("call back after ping err: ", err)
//
//	}
//}

func main() {
	s := hnet.NewServer("[hinx v01]")

	s.SetOnConnStart(func(connection hiface.IConnection) {
		fmt.Println("==> DoConnectionBegin is Called...")
		connection.SetProperty("Server", "hinx")
	})
	s.SetOnConnStop(func(connection hiface.IConnection) {
		fmt.Println("==> DoConnectionClose is Called...")

		if server, err := connection.GetProperty("Server"); err == nil {
			fmt.Println("Server = ", server)
		}
	})

	s.AddRouter(0, &PingRouter{})

	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
