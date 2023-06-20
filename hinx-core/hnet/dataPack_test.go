package hnet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7788")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go func(conn net.Conn) {
			dp := NewDataPack()
			for {
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println(err)
					return
				}

				unpack, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println(err)
					return
				}

				if unpack.GetDataLen() > 0 {
					msg := unpack.(*Message)
					msg.Data = make([]byte, msg.GetDataLen())

					_, err = io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println(err)
						return
					}

					fmt.Println("Recv msgId: ", msg.Id, "dataLen: ", msg.DataLen, "data: ", string(msg.Data))
				}
			}
		}(conn)
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:7788")
	if err != nil {
		fmt.Println(err)
		return
	}

	dp := NewDataPack()

	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte("hello"),
	}

	msg2 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte("world"),
	}
	pack1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println(err)
		return
	}
	pack2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.Write(append(pack1, pack2...))

	time.Sleep(5 * time.Second)
}
