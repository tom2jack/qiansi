package deploy

import "testing"

func TestRunShell(t *testing.T) {
	RunShell("cs:\\", `
piqng baidu.com

`)
	RunShell("c:\\", `
pinqg baidu.com
`)

}
