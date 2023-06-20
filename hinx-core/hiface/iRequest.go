package hiface

// IRequest compact connection info and request data to a request
type IRequest interface {
	// GetConnection get connection
	GetConnection() IConnection
	// GetData get data
	GetData() []byte
	GetMsgID() uint32
}
