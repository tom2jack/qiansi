package resp

import "time"

type ServerVO struct {
	CreateTime   time.Time
	DeviceId     string
	Domain       string
	Id           int
	ServerName   string
	ServerRuleId int
	ServerStatus int
	UpdateTime   time.Time
}
