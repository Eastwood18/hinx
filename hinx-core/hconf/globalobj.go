package hconf

import (
	"gopkg.in/yaml.v3"
	"hinx/hinx-core/hiface"
	"hinx/utils/commandLine/args"
	"hinx/utils/commandLine/uflag"
	"os"
	"testing"
)

var GlobalObject *GlobalObj

type GlobalObj struct {
	/*
		Server
	*/
	Server  hiface.IServer
	Host    string
	TcpPort int
	Name    string

	/*
		hinx
	*/
	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32

	Heartbeat int
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func (g *GlobalObj) Reload() {

	////fmt.Println(os.Getwd())
	//file, err := os.ReadFile("conf/hinx.yaml")
	//if err != nil {
	//	panic(err)
	//}
	//err = yaml.Unmarshal(file, &GlobalObject)
	//if err != nil {
	//	panic(err)
	//}

	confFilePath := args.Args.ConfigFile
	if confFileExists, _ := PathExists(confFilePath); confFileExists != true {

		// The configuration file may not exist,
		// in which case the default parameters should be used to initialize the logging module configuration.
		// (配置文件不存在也需要用默认参数初始化日志模块配置)
		//g.InitLogConfig()
		//
		//zlog.Ins().ErrorF("Config File %s is not exist!!", confFilePath)
		return
	}

	data, err := os.ReadFile("conf/hinx.yaml")
	if err != nil {
		panic(err)
	}

	yaml.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}

}

/*
	define a globalobj
*/

func init() {

	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	args.InitConfigFlag(pwd+"/conf/hinx.yaml", "The configuration file defaults to <exeDir>/conf/zinx.json if it is not set.")
	testing.Init()
	uflag.Parse()

	// after parsing
	args.FlagHandle()
	GlobalObject = &GlobalObj{
		Server:           nil,
		Host:             "0.0.0.0",
		TcpPort:          8899,
		Name:             "HinxServerApp",
		Version:          "0.1",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   8,
		MaxWorkerTaskLen: 1024,
		Heartbeat:        30000,
	}

	// reload conf/hinx.yaml
	GlobalObject.Reload()
	//fmt.Println(GlobalObject)
}
