package utils

import (
	"errors"
	"strings"
)

func ParseQiansiClientFileInfo(fileName string) (version, os, arch string, err error) {
	splitName := strings.Split(fileName, "_")
	if len(splitName) != 3 {
		err = errors.New("can't parse the name")
		return
	}
	version = splitName[1]
	splitArch := strings.Split(splitName[2], "-")
	if len(splitArch) != 2 {
		err = errors.New("can't parse the arch")
		return
	}
	os, arch = splitArch[0], splitArch[1]
	return
}
