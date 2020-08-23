package models

type mqttModels struct{}

func GetMQTTModels() *mqttModels {
	return &mqttModels{}
}

// TODO

//CreateUser 创建Mqtt用户
func (m *mqttModels) CreateUser(username, password, isSuper bool) error {
	return nil
}

//DeleteUserByName 根据用户名删除用户
func (m *mqttModels) DeleteUserByName(userName string) error {
	return nil
}

//RestPassword 重置密码
func (m *mqttModels) RestPassword(userName, newPwd string) error {
	return nil
}

//AddClientAcl 添加客户端权限
func (m *mqttModels) AddClientAcl(username string) error {
	return nil
}
