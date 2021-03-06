package models

import "time"

// Deploy 纸喵部署应用表
type Deploy struct {
	ID            int       `gorm:"primary_key;column:id;type:int(10) unsigned;not null" json:"id"`          // 应用唯一编号
	UId           int       `gorm:"index:IX_UID;column:uid;type:int(11);not null" json:"uid"`                // 创建用户UID
	Title         string    `gorm:"column:title;type:varchar(120);not null" json:"title"`                    // 应用名称
	DeployType    int8      `gorm:"column:deploy_type;type:tinyint(1);not null" json:"deploy_type"`          // 部署类型 0-本地 1-git 2-zip 3-docker
	WorkDir       string    `gorm:"column:work_dir;type:varchar(255);not null" json:"work_dir"`              // 工作目录
	BeforeCommand string    `gorm:"column:before_command;type:varchar(2000);not null" json:"before_command"` // 前置命令
	AfterCommand  string    `gorm:"column:after_command;type:varchar(2000);not null" json:"after_command"`   // 后置命令
	NowVersion    int       `gorm:"column:now_version;type:int(10) unsigned;not null" json:"now_version"`    // 当前版本
	OpenID        string    `gorm:"unique;column:open_id;type:char(32);not null" json:"open_id"`             // 应用开放编码(用于hook部署)
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

// DeployDocker 纸喵部署-docker
type DeployDocker struct {
	DeployID         int       `gorm:"primary_key;column:deploy_id;type:int(10) unsigned;not null" json:"deploy_id"`  // 部署应用编号
	DockerImage      string    `gorm:"column:docker_image;type:varchar(500);not null" json:"docker_image"`            // 资源地址(完整路径)
	UserName         string    `gorm:"column:user_name;type:varchar(255);not null" json:"user_name"`                  // 账号
	Password         string    `gorm:"column:password;type:varchar(255);not null" json:"password"`                    // 密码
	IsRuning         int8      `gorm:"column:is_runing;type:tinyint(1);not null" json:"is_runing"`                    // 是否启动 -1-不启动 1-启动
	ContainerName    string    `gorm:"column:container_name;type:varchar(255);not null" json:"container_name"`        // 容器名称
	ContainerVolumes string    `gorm:"column:container_volumes;type:varchar(1000);not null" json:"container_volumes"` // 地址映射 a:b a宿主机 b容器内
	ContainerPorts   string    `gorm:"column:container_ports;type:varchar(1000);not null" json:"container_ports"`     // 端口暴露 a:b a内部 b外部
	ContainerEnv     string    `gorm:"column:container_env;type:varchar(2000);not null" json:"container_env"`         // 环境变量注入
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

// DeployGit 纸喵部署-git
type DeployGit struct {
	DeployID   int       `gorm:"primary_key;column:deploy_id;type:int(10) unsigned;not null" json:"deploy_id"` // 部署应用编号
	RemoteURL  string    `gorm:"column:remote_url;type:varchar(500);not null" json:"remote_url"`               // 资源地址
	DeployPath string    `gorm:"column:deploy_path;type:varchar(500);not null" json:"deploy_path"`             // 本地部署地址
	Branch     string    `gorm:"column:branch;type:varchar(100);not null" json:"branch"`                       // git分支
	DeployKeys string    `gorm:"column:deploy_keys;type:varchar(3000);not null" json:"deploy_keys"`            // 部署私钥
	UserName   string    `gorm:"column:user_name;type:varchar(255);not null" json:"user_name"`                 // 账号
	Password   string    `gorm:"column:password;type:varchar(255);not null" json:"password"`                   // 密码
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

// DeployServerRelation 部署服务-server关联表
type DeployServerRelation struct {
	DeployID      int       `gorm:"primary_key;column:deploy_id;type:int(11);not null" json:"deploy_id"` // 部署服务id
	ServerID      int       `gorm:"primary_key;column:server_id;type:int(11);not null" json:"server_id"` // 服务器ID
	DeployVersion int       `gorm:"column:deploy_version;type:int(10) unsigned" json:"deploy_version"`   // 已部署版本
	CreateTime    time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

// DeployZip 纸喵部署-zip
type DeployZip struct {
	DeployID   int       `gorm:"primary_key;column:deploy_id;type:int(10) unsigned;not null" json:"deploy_id"` // 部署应用编号
	RemoteURL  string    `gorm:"column:remote_url;type:varchar(500);not null" json:"remote_url"`               // 资源地址
	DeployPath string    `gorm:"column:deploy_path;type:varchar(500);not null" json:"deploy_path"`             // 本地部署地址
	Password   string    `gorm:"column:password;type:varchar(255);not null" json:"password"`                   // 密码
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}
