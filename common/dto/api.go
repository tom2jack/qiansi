package dto

import "time"

type DeployDTO struct {
	AfterCommand  string    `xorm:"not null comment('后置命令') VARCHAR(2000)"`
	BeforeCommand string    `xorm:"not null default '' comment('前置命令') VARCHAR(2000)"`
	Branch        string    `xorm:"not null default '' comment('git分支') VARCHAR(100)"`
	CreateTime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	DeployKeys    string    `xorm:"not null default '' comment('部署私钥') VARCHAR(3000)"`
	DeployType    int       `xorm:"not null default 0 comment('部署类型 0-本地 1-git 2-zip') TINYINT(1)"`
	Id            int       `xorm:"not null pk autoincr comment('应用唯一编号') index(IX_UID) INT(10)"`
	LocalPath     string    `xorm:"not null default '' comment('本地部署地址') VARCHAR(500)"`
	NowVersion    int       `xorm:"not null default 0 comment('当前版本') INT(10)"`
	RemoteUrl     string    `xorm:"not null default '' comment('资源地址') VARCHAR(500)"`
	Title         string    `xorm:"not null default '' comment('应用名称') VARCHAR(120)"`
	Uid           int       `xorm:"not null comment('创建用户UID') index(IX_UID) INT(11)"`
	UpdateTime    time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}