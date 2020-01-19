package deploy

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common/dto"
	"gitee.com/zhimiao/qiansi/common/logger"
	"gitee.com/zhimiao/qiansi/common/utils"
	"gitee.com/zhimiao/qiansi/qiansi-client/request"
	"github.com/progrium/go-shell"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	Task *utils.SafeStringMap
	RUN  bool = false
)

func init() {
	if RUN {
		return
	}
	Task = utils.NewSafeStringMap()
	RUN = true
}

func Run(data string) {
	TaskList := []dto.DeployDTO{}
	// 此处请求只有一次机会
	_ = request.GetDeployTask(&TaskList)
	for _, v := range TaskList {
		go runTask(&v)
	}
}

func runTask(deploy *dto.DeployDTO) {
	logger.Info("%d号部署任务正在启动...", deploy.Id)
	if !Task.SETNX(strconv.Itoa(deploy.Id), strconv.FormatInt(time.Now().Unix(), 36)) {
		logger.Warn("%d号部署任务未完成，拒绝执行", deploy.Id)
		return
	}
	defer Task.DEL(strconv.Itoa(deploy.Id))
	// 执行前置命令
	if deploy.BeforeCommand != "" {
		RunShell(deploy, deploy.BeforeCommand)
	}
	// 选择部署操作
	switch deploy.DeployType {
	case 1:
		Git(deploy)
	}
	// 执行后置命令
	if deploy.AfterCommand != "" {
		RunShell(deploy, deploy.AfterCommand)
	}
	request.DeployNotify(deploy)
}

func LogPush(deploy *dto.DeployDTO, format string, v ...interface{}) {
	var fname = "default"
	if pc, _, _, ok := runtime.Caller(1); ok {
		fname = runtime.FuncForPC(pc).Name()
	}
	logger.Info("("+fname+") "+format, v...)
	// fname = strings.ReplaceAll(fname, "gitee.com/zhimiao/qiansi/qiansi-client/deploy.", "")
	// 反向推送日志到千丝平台
	request.LogPush(deploy, fmt.Sprintf("("+fname+") "+format, v...))
}

func RunShell(deploy *dto.DeployDTO, cmd string) error {
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
		res.SetWorkDir(deploy.LocalPath)
		o := res.Run()
		o.Wait()
		out := o.String()
		if !utf8.ValidString(out) {
			out = utils.ConvertToString(out, "gbk", "utf8")
		}
		LogPush(deploy, "正在执行命令: %s, \n结果输出:\n %s\n错误输出:\n %s", v, out, o.Error().Error())
	}
	return nil
}
