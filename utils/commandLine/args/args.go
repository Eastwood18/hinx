package args

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hinx/utils/commandLine/uflag"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

type args struct {
	ExeAbsDir   string
	ExeName     string
	ConfigFile  string
	MachineCode uint
}

var (
	Args   = args{}
	isInit = false
)

func init() {
	exe := os.Args[0]

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	Args.ExeAbsDir = pwd
	Args.ExeName = filepath.Base(exe)
	Args.MachineCode = getMachineCode()
}

func InitConfigFlag(defaultValue string, tips string) {
	if isInit {
		return
	}
	isInit = true

	uflag.StringVar(&Args.ConfigFile, "c", defaultValue, tips)
	return
}

func FlagHandle() {
	filePath, err := filepath.Abs(Args.ConfigFile)
	if err != nil {
		panic(err)
	}
	Args.ConfigFile = filePath
}

func getMachineCode() uint {
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorID")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	//	fmt.Println(string(out))
	str := string(out)
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	str = reg.ReplaceAllString(str, "")

	hashbyte := md5.Sum([]byte(str[11:]))
	hashint := uint(binary.BigEndian.Uint32(hashbyte[:]))

	return hashint
}
