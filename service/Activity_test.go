package service

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestInviterActive(t *testing.T) {
	for i := 20; i > 0; i++ {
		x := rand.Intn(2)
		fmt.Println(x)
	}
}
