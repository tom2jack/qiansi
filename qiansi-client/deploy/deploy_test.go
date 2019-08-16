package deploy

import "testing"

func TestRunShell(t *testing.T) {
	_ = RunShell("cs:\\", `
piqng baidu.com

`)
	_ = RunShell("c:\\", `
ping baidu.com
`)

}
