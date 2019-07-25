package deploy

import (
	"fmt"
	"github.com/progrium/go-shell"
	"qiansi/common/models"
	"qiansi/common/utils"
	"qiansi/common/zmlog"
	"qiansi/qiansi-client/request"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	Task *utils.SafeMap
	RUN  bool = false
)

func init() {
	if RUN {
		return
	}
	Task = utils.NewSafeMap()
	RUN = true
}

func Run() {
	// TODO: 原子执行逻辑
	TaskList := []models.Deploy{}
	_ = request.GetDeployTask(&TaskList)
	for _, v := range TaskList {
		go runTask(v)
	}
}

func runTask(deploy models.Deploy) {
	if !Task.SETNX(strconv.Itoa(deploy.Id), strconv.FormatInt(time.Now().Unix(), 36)) {
		return
	}
	// 执行前置命令
	if deploy.BeforeCommand != "" {
		RunShell(deploy.LocalPath, deploy.BeforeCommand)
	}
	// 选择部署操作
	switch deploy.DeployType {
	case 1:
		Git(&deploy)
	}
	// 执行后置命令
	if deploy.AfterCommand != "" {
		RunShell(deploy.LocalPath, deploy.AfterCommand)
	}
	Task.DEL(strconv.Itoa(deploy.Id))
}

func LogPush(format string, v ...interface{}) {
	// TODO: 推送至千丝平台
	var fname = "default"
	if pc, _, _, ok := runtime.Caller(1); ok {
		fname = runtime.FuncForPC(pc).Name()
	}
	zmlog.Info("("+fname+") "+format, v...)
	// 反向推送日志到千丝平台
	request.LogPush(fmt.Sprintf("("+fname+") "+format, v...))
}

func RunShell(work_path, cmd string) error {
	defer shell.ErrExit()
	if runtime.GOOS == "windows" {
		shell.Shell = []string{"cmd", "/c"}
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
