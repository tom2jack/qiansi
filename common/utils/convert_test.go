package utils

import (
	"fmt"
	"testing"
)

type B struct {
	A1 string
	B1 int
	G1 bool
}

type A struct {
	A1 string
}

func TestSuperConvert(t *testing.T) {
	a := A{
		A1: "123",
	}
	b := B{}
	SuperConvert(&a, &b)
	fmt.Printf("%#v", b)
}
