package common

import (
	"testing"
)

func TestSave(t *testing.T) {
	Config.Init()
	Config.Save()
}
