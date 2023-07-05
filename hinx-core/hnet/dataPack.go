package hnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hinx/hinx-core/hconf"
	"hinx/hinx-core/hiface"
)

// DataPack base TLV format unpack message
type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}
func (d *DataPack) GetHeadLen() uint32 {
	return 8 // DataLen + ID
}

func (d *DataPack) Pack(message hiface.IMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.LittleEndian, message.GetDataLen())
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, message.GetMsgId())
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, message.GetMsgData())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *DataPack) Unpack(binaryData []byte) (hiface.IMessage, error) {
	buf := bytes.NewReader(binaryData)
	msg := new(Message)
	err := binary.Read(buf, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}
	if hconf.GlobalObject.MaxPackageSize > 0 && msg.DataLen > hconf.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg recv")
	}
	return msg, nil
}
