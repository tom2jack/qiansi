package deploy

import "testing"

func TestRunShell(t *testing.T) {
	RunShell("C:\\", `
dir

tasklist

`)
}
