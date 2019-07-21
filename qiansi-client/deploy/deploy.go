package deploy

import (
	"github.com/progrium/go-shell"
	"qiansi/common/models"
	"qiansi/common/zmlog"
	"qiansi/qiansi-client/request"
	"runtime"
	"strings"
)

func Run() {
	TaskList := []models.Deploy{}
	_ = request.GetDeployTask(&TaskList)
	for _, v := range TaskList {
		// 执行前置命令
		if v.BeforeCommand != "" {
			RunShell(v.LocalPath, v.BeforeCommand)
		}
		// 选择部署操作
		switch v.DeployType {
		case 1:
			Git(&v)
		}
		// 执行后置命令
		if v.AfterCommand != "" {
			RunShell(v.LocalPath, v.AfterCommand)
		}
	}
}

func LogPush(format string, v ...interface{}) {
	// TODO: 推送至千丝平台
	var fname = "default"
	if pc, _, _, ok := runtime.Caller(1); ok {
		fname = runtime.FuncForPC(pc).Name()
	}
	zmlog.Info("("+fname+") "+format, v...)
}

func RunShell(work_path, cmd string) error {
	defer shell.ErrExit()
	if runtime.GOOS == "windows" {
		shell.Shell = []string{"C:\\Windows\\System32\\cmd.exe", "/c"}
	}
	shell.Panic = false
	cmd = strings.ReplaceAll(cmd, "\r\n", "\n")
	cmds := strings.Split(cmd, "\n")
	for _, v := range cmds {
		if v == "" {
			continue
		}
		res := shell.Cmd(v)
		res.SetWorkDir(work_path)
		o := res.Run()
		o.Wait()
		LogPush("正在执行命令: %s, \n结果输出:\n %s\n错误输出:\n %s", v, o.String(), o.Error().Error())
	}
	return nil
}
