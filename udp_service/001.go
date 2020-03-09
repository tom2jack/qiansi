package udp_service

import (
	"bytes"
	"encoding/gob"
	"gitee.com/zhimiao/qiansi/dto"
	"strconv"
	"sync"
)

// 001-轮训任务存储器
type Hook001PO struct {
	sync.RWMutex
	// 部署存储集合
	Deploy map[string]string
	// 调度任务
	Scheduld map[string]string
	// 指标监控配置
	Telegraf map[string]string
}

var Hook001 *Hook001PO

func init() {
	Hook001 = &Hook001PO{
		Deploy:   make(map[string]string),
		Scheduld: make(map[string]string),
		Telegraf: make(map[string]string),
	}
}

// GetHookData 获取存储器交换数据
func (m *Hook001PO) GetHookData(key []byte) dto.Hook001DTO {
	skey := string(key)
	m.RLock()
	dtoData := dto.Hook001DTO{
		Deploy:   m.Deploy[skey],
		Scheduld: m.Scheduld[skey],
		Telegraf: m.Telegraf[skey],
	}
	m.RUnlock()
	return dtoData
}

// AddDeploy 添加部署消息
func (m *Hook001PO) AddDeploy(serverId int) {
	s := strconv.Itoa(serverId)
	if s != "" {
		m.Lock()
		m.Deploy[s] = "1"
		m.Unlock()
	}
}

// DelDeploy 删除部署消息
func (m *Hook001PO) DelDeploy(serverId int) {
	s := strconv.Itoa(serverId)
	if s != "" {
		m.Lock()
		delete(m.Deploy, s)
		m.Unlock()
	}
}

// AddScheduld 添加调度消息
func (m *Hook001PO) AddScheduld(serverId int) {
	s := strconv.Itoa(serverId)
	if s != "" {
		m.Lock()
		m.Scheduld[s] = "1"
		m.Unlock()
	}
}

// DelScheduld 删除调度消息
func (m *Hook001PO) DelScheduld(serverId int) {
	s := strconv.Itoa(serverId)
	if s != "" {
		m.Lock()
		delete(m.Scheduld, s)
		m.Unlock()
	}
}

// AddTelegraf 添加指标消息
func (m *Hook001PO) AddTelegraf(serverId int) {
	s := strconv.Itoa(serverId)
	if s != "" {
		m.Lock()
		m.Telegraf[s] = "1"
		m.Unlock()
	}
}

// DelTelegraf 删除指标消息
func (m *Hook001PO) DelTelegraf(serverId int) {
	s := strconv.Itoa(serverId)
	if s != "" {
		m.Lock()
		delete(m.Telegraf, s)
		m.Unlock()
	}
}

// ClientTaskLoop 轮训， Task结构为 {001{deploy_id}:1}
func Hook_001(request []byte) []byte {
	var network bytes.Buffer
	// 1.创建编码器
	enc := gob.NewEncoder(&network)
	// 2.向编码器中写入数据
	err := enc.Encode(Hook001.GetHookData(request))
	if err != nil {
		return nil
	}
	return network.Bytes()
}
