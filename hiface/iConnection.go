package hiface

import "net"

// IConnection define abstract connection module
type IConnection interface {
	// Start  connection
	Start()
	// Stop  connection
	Stop()
	// GetTCPConnection get connection socket conn
	GetTCPConnection() *net.TCPConn
	// GetConnID get connection connect_id
	GetConnID() uint32
	// GetRemoteAddr get remote client IP Port
	GetRemoteAddr() net.Addr
	// SendMsg Send  message
	SendMsg(msgId uint32, data []byte) error

	// SetProperty set connection property
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}

// HandleFunc define a function to handle connection business
type HandleFunc func(*net.TCPConn, []byte, int) error
