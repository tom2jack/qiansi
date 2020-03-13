package notifyevent

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"sync"
	"time"
)

// 001-轮训任务存储器
type Hook001PO struct {
	sync.RWMutex
	// 部署存储集合
	Deploy map[int]string
	// 调度任务
	Scheduld map[int]string
	// 指标监控配置
	Telegraf map[int]string
	// 活跃服务器
	ActiveServer map[int]time.Time
}

// 001-轮训任务数据交换
type Hook001DTO struct {
	Deploy   string
	Scheduld string
	Telegraf string
}

var Hook001 *Hook001PO

func init() {
	Hook001 = &Hook001PO{
		Deploy:       make(map[int]string),
		Scheduld:     make(map[int]string),
		Telegraf:     make(map[int]string),
		ActiveServer: make(map[int]time.Time),
	}
}

// GetHookData 获取存储器交换数据
func (m *Hook001PO) GetHookData(serverId int) Hook001DTO {
	m.RLock()
	dtoData := Hook001DTO{
		Deploy:   m.Deploy[serverId],
		Scheduld: m.Scheduld[serverId],
		Telegraf: m.Telegraf[serverId],
	}
	m.RUnlock()
	return dtoData
}

// AddDeploy 添加部署消息
func (m *Hook001PO) AddDeploy(serverId int) {
	m.Lock()
	m.Deploy[serverId] = "1"
	m.Unlock()
}

// DelDeploy 删除部署消息
func (m *Hook001PO) DelDeploy(serverId int) {
	m.Lock()
	delete(m.Deploy, serverId)
	m.Unlock()
}

// AddScheduld 添加调度消息
func (m *Hook001PO) AddScheduld(serverId int) {
	m.Lock()
	m.Scheduld[serverId] = "1"
	m.Unlock()
}

// DelScheduld 删除调度消息
func (m *Hook001PO) DelScheduld(serverId int) {
	delete(m.Scheduld, serverId)
}

// AddTelegraf 添加指标消息
func (m *Hook001PO) AddTelegraf(serverId int) {
	m.Lock()
	m.Telegraf[serverId] = "1"
	m.Unlock()
}

// DelTelegraf 删除指标消息
func (m *Hook001PO) DelTelegraf(serverId int) {
	m.Lock()
	delete(m.Telegraf, serverId)
	m.Unlock()
}

func (m *Hook001PO) UpdateServerLastRequest(serverId int) {
	m.Lock()
	m.ActiveServer[serverId] = time.Now()
	m.Unlock()
}

// GetServerLastRequest 获取服务器最后一次请求
func (m *Hook001PO) GetServerLastRequest(serverId ...int) (r map[int]time.Time) {
	m.RLock()
	for _, v := range serverId {
		if t, ok := m.ActiveServer[v]; ok {
			r[v] = t
		}
	}
	m.RUnlock()
	return
}

// GetServerLastRequest 获取服务器最后一次请求
func (m *Hook001PO) GetActiveServerNum(serverId ...int) (num int) {
	ts := time.Now().Add(-10 * time.Second)
	m.RLock()
	for _, v := range serverId {
		if t, ok := m.ActiveServer[v]; ok && t.After(ts) {
			num++
		}
	}
	m.RUnlock()
	return
}

// ClientTaskLoop 轮训
func Hook_001(request []byte) []byte {
	serverId, err := strconv.Atoi(string(request))
	if err != nil {
		return nil
	}
	// 更新最后请求
	Hook001.UpdateServerLastRequest(serverId)
	var network bytes.Buffer
	// 1.创建编码器
	enc := gob.NewEncoder(&network)
	// 2.向编码器中写入数据
	err = enc.Encode(Hook001.GetHookData(serverId))
	if err != nil {
		return nil
	}
	return network.Bytes()
}
