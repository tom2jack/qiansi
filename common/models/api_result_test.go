package models

import "testing"

func TestNewApiResult(t *testing.T) {
	t.Log(*NewApiResult(8, nil, "成功"))
}
