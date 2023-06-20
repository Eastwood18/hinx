package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"hinx/hinx-core/hiface"
	"os"
)

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer hiface.IServer
	Host      string
	TcpPort   int
	Name      string

	/*
		hinx
	*/
	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
}

func (g *GlobalObj) Reload() {

	fmt.Println(os.Getwd())
	file, err := os.ReadFile("conf/hinx.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	define a globalobj
*/

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		TcpServer:        nil,
		Host:             "0.0.0.0",
		TcpPort:          8899,
		Name:             "HinxServerApp",
		Version:          "0.1",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   8,
		MaxWorkerTaskLen: 1024,
	}

	// reload conf/hinx.yaml
	GlobalObject.Reload()
	fmt.Println(GlobalObject)
}
