package resp

type ServerVO struct {
	ApiSecret    string
	CreateTime   JsonTimeDate
	DeviceId     string
	Domain       string
	Id           int
	ServerName   string
	ServerRuleId int
	ServerStatus int
	Uid          int
	UpdateTime   JsonTimeDate
}
