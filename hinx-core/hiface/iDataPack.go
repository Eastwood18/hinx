package hiface

/*
	pack, unpack module
	orient to TCP data stream
*/

// IDataPack base TLV format pack message
type IDataPack interface {
	GetHeadLen() uint32
	Pack(message IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
