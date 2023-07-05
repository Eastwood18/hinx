package main

import (
	"fmt"
	"hinx/hinx-core/hnet"
	"io"
	"net"
	"strconv"
	"time"
)

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8899")

	for i := 0; i < 10; i++ {
		dp := hnet.NewDataPack()
		pack, err := dp.Pack(hnet.NewMsgPackage(1, []byte("hello hinx")))
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, v := range pack {
			fmt.Print(strconv.FormatInt(int64(v), 16) + " ")
		}
		conn.Write(pack)

		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println(err)
			return
		}

		unpack, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(unpack)
		if unpack.GetDataLen() > 0 {
			msg := unpack.(*hnet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err = io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Recv msgId: ", msg.Id, "dataLen: ", msg.DataLen, "data: ", string(msg.Data))

		}
		time.Sleep(5 * time.Second)
	}

}
