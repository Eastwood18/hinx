package utils

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"os/exec"
	"regexp"
	"testing"
)

func TestSnowFlake(t *testing.T) {

	go func() {
		for i := 0; i < 10; i++ {
			id, _ := GetSnowFlakeID()
			fmt.Printf("%b\n", id)
		}
	}()
}

func TestGetCpuId(t *testing.T) {
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
	fmt.Println(hashint)
	fmt.Println(str)
}
