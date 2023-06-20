package hiface

// IServer define a server interface
type IServer interface {
	// Start  server
	Start()
	// Stop server
	Stop()
	// Serve server
	Serve()

	// router function
	AddRouter(id uint32, router IRouter)

	GetConnManager() IConnManager

	SetOnConnStart(func(connection IConnection))
	SetOnConnStop(func(connection IConnection))

	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)
}
