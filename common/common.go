package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strings"
)

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

//GetLocalMac 随便获取本机的一个Mac地址
func GetLocalMac() string {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces {
		if inter.HardwareAddr.String() != "" {
			return inter.HardwareAddr.String()
		}
	}
	return ""
}

//GetRandStr 获取随机字符串
func GetRandStr(len int) string {
	b := make([]byte, len)
	rand.Read(b) //在byte切片中随机写入元素
	return base64.StdEncoding.EncodeToString(b)
}

func DecrptogAES(src, key string) string {
	s := []byte(src)
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return ""
	}
	blockMode := cipher.NewCBCDecrypter(block, k)
	blockMode.CryptBlocks(s, s)
	n := len(s)
	count := int(s[n-1])
	return string(s[:n-count])
}
