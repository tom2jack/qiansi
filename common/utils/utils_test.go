package utils

import (
	"fmt"
	"testing"
)

func TestParseQiansiClientFileInfo(t *testing.T) {
	info, os, arch, err := ParseQiansiClientFileInfo("qiansi-client_1.3.1_linux-amd64")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(info)
	fmt.Println(os)
	fmt.Println(arch)
}
