package hnet

import (
	hiface2 "hinx/hinx-core/hiface"
)

type Request struct {
	conn hiface2.IConnection //  connection established to client
	msg  hiface2.IMessage
}

func (r *Request) GetConnection() hiface2.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
