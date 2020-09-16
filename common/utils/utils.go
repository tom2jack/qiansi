package utils

import (
	"errors"
	"strings"
)

// ParseQiansiClientFileInfo 解析千丝客户端文件名
// qiansi-client_1.0.0_windows-amd64.exe
func ParseQiansiClientFileInfo(fileName string) (version, os, arch string, err error) {
	splitName := strings.Split(fileName, "_")
	nlen := len(splitName)
	if nlen < 3 {
		err = errors.New("can't parse the name")
		return
	}
	version = splitName[nlen-2]
	splitArch := strings.Split(splitName[nlen-1], "-")
	if len(splitArch) != 2 {
		err = errors.New("can't parse the arch")
		return
	}
	os = splitArch[0]
	arch = splitArch[1]
	extIndex := strings.Index(arch, ".")
	if extIndex > 0 {
		arch = arch[:extIndex]
	}
	return
}
