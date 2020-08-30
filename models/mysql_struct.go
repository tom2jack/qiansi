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

// MqttACL [...]
type MqttACL struct {
	ID       int    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	Allow    int    `gorm:"column:allow;type:int(1)" json:"allow"`                // 0: deny, 1: allow
	IPaddr   string `gorm:"column:ipaddr;type:varchar(60)" json:"ipaddr"`         // IpAddress
	Username string `gorm:"column:username;type:varchar(100)" json:"username"`    // Username
	Clientid string `gorm:"column:clientid;type:varchar(100)" json:"clientid"`    // ClientId
	Access   int    `gorm:"column:access;type:int(2);not null" json:"access"`     // 1: subscribe, 2: publish, 3: pubsub
	Topic    string `gorm:"column:topic;type:varchar(100);not null" json:"topic"` // Topic Filter
}

// MqttUser [...]
type MqttUser struct {
	ID          int       `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	Username    string    `gorm:"unique;column:username;type:varchar(100)" json:"username"`
	Password    string    `gorm:"column:password;type:varchar(100)" json:"password"`
	Salt        string    `gorm:"column:salt;type:varchar(35)" json:"salt"`
	IsSuperuser bool      `gorm:"column:is_superuser;type:tinyint(1)" json:"is_superuser"`
	Created     time.Time `gorm:"column:created;type:datetime" json:"created"`
}

// Member 用户表
type Member struct {
	Id          int       `gorm:"primary_key;column:id;type:int(10) unsigned;not null"`
	Phone       string    `gorm:"unique;column:phone;type:char(11);not null"`         // 手机号
	Password    string    `gorm:"column:password;type:varchar(255);not null"`         // 密码
	InviterUid  int       `gorm:"column:inviter_uid;type:int(10) unsigned;not null"`  // 邀请人UID
	MaxDeploy   int       `gorm:"column:max_deploy;type:int(10) unsigned;not null"`   // 最大部署应用数量
	MaxSchedule int       `gorm:"column:max_schedule;type:int(10) unsigned;not null"` // 最大调度任务数量
	Status      int8      `gorm:"column:status;type:tinyint(1);not null"`             // 0-锁定 1-正常
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;not null"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;not null"`
}

// Server 服务器注册表
type Server struct {
	ID           int       `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	UId          int       `gorm:"index:IX_uid;column:uid;type:int(10) unsigned;not null" json:"uid"`  // 用户ID
	ServerName   string    `gorm:"column:server_name;type:varchar(64);not null" json:"server_name"`    // 服务器备注名
	ServerStatus int8      `gorm:"column:server_status;type:tinyint(3);not null" json:"server_status"` // 服务器状态 -1-失效 0-待认领 1-已分配通信密钥 2-已绑定
	MqttUser     string    `gorm:"column:mqtt_user;type:varchar(255);not null" json:"mqtt_user"`       // mqtt用户名
	APISecret    string    `gorm:"column:api_secret;type:varchar(32);not null" json:"api_secret"`      // API密钥
	DeviceID     string    `gorm:"column:device_id;type:char(36);not null" json:"device_id"`           // 服务器唯一设备号
	Domain       string    `gorm:"column:domain;type:varchar(255);not null" json:"domain"`             // 服务器地址(域名/ip)
	CreateTime   time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}
