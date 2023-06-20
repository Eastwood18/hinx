package hiface

// IMessage packing request message
type IMessage interface {
	GetMsgId() uint32
	GetDataLen() uint32
	GetMsgData() []byte

	SetMsgId(uint32)
	SetDataLen(uint32)
	SetMsgData([]byte)
}
