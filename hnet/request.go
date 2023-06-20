package hnet

import "hinx/hiface"

type Request struct {
	conn hiface.IConnection //  connection established to client
	msg  hiface.IMessage
}

func (r *Request) GetConnection() hiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
