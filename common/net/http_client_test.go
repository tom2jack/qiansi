package net

import (
	"fmt"
	"testing"
	"time"
)

func TestHttpClient_Get(t *testing.T) {
	for i := 0; i < 400000; i++ {
		go func(i int) {
			for x := 0; x < 20; x++ {
				time.Sleep(time.Second)
			}

			// r := HttpClient.Get("http://192.168.1.24:1301/log.php?x=123", 30)
			// fmt.Printf("%d: %#v\n", i, r)
			fmt.Printf("%d--结束\n", i)
		}(i)
	}
	fmt.Printf("jieshu")
	select {}
}
