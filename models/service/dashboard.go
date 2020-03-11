package service

import "gitee.com/zhimiao/qiansi/models"

type MaxInfo struct {
	DeployNum   int
	MaxDeploy   int
	ScheduleNum int
	MaxSchedule int
}

// GetUserModuleMaxInfo 获取用户模块限额信息
func GetUserModuleMaxInfo(uid int) (info *MaxInfo, err error) {
	member := models.Member{Id: uid}
	err = member.MaxInfo()
	if err != nil {
		return
	}
	info.MaxDeploy, info.MaxSchedule = member.MaxDeploy, member.MaxSchedule
	var num int
	deploy := models.Deploy{Uid: uid}
	num, err = deploy.Count()
	if err != nil {
		return
	}
	info.DeployNum = num
	schedule := models.Schedule{Uid: uid}
	num, err = schedule.Count()
	if err != nil {
		return
	}
	info.ScheduleNum = num
	return
}
