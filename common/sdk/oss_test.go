package sdk

import (
	"fmt"
	"github.com/zhi-miao/qiansi/common/config"
	"testing"
)

func TestNewOSSClient(t *testing.T) {
	config.LoadConfig("config.toml")
	client, err := NewOSSClient()
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err := client.ListFile("qiansi-client")
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, s := range file {

		fmt.Println(s)
	}
}
