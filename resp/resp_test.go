package resp

import (
	"fmt"
	"testing"
)

func TestNewApiResult(t *testing.T) {
	err := fmt.Errorf("sdfsdfasdf")
	// err = nil
	d := NewApiResult(err)
	fmt.Print(d)
}
