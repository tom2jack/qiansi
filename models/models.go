package models

import (
	"context"
	"fmt"
	"github.com/zhi-miao/qiansi/common/config"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"github.com/influxdata/influxdb-client-go/api/write"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Redis    *redis.Client
	Mysql    *gorm.DB
	InfluxDB *zmInflux
)

const (
	UserServerIdsCacheKey = "QIANSI:dashboard:user-server-ids:%d"
	// 设备上下线状态
	ServerOnlineStatusCacheKey = "QIANSI:ServerOnlineStatus"
	// 设备号转服务器编号
	ServerDeviceIDCacheKey = "QIANSI:ServerDeviceID"
)

type zmInflux struct {
	sync.Mutex
	Client        influxdb2.Client
	WriteApiCache map[string]api.WriteApi
	ReadApiCache  map[string]api.QueryApi
}

type CommonMap map[string]interface{}

// TempModelStruct 临时数据结构体
type TempModelStruct struct {
	Num int  `gorm:"column:num"`
	Has bool `gorm:"column:has"`
}

// Start 初始化数据
func Start() {
	loadRedis()
	loadMysql()
	loadInfluxDB()
}

// 初始化influxDb
func loadInfluxDB() {
	InfluxDB = &zmInflux{
		Client:        influxdb2.NewClient(config.GetConfig().InfluxDB.Host, config.GetConfig().InfluxDB.Token),
		WriteApiCache: make(map[string]api.WriteApi),
		ReadApiCache:  make(map[string]api.QueryApi),
	}
}

// Setup Initialize the Redis instance
func loadRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.GetConfig().Redis.Host,
		Password: config.GetConfig().Redis.Auth, // no password set
		DB:       config.GetConfig().Redis.DB,   // use default DB
	})
}

// Setup Initialize the Mysql instance
func loadMysql() {
	var err error
	Mysql, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.GetConfig().Mysql.User,
		config.GetConfig().Mysql.Password,
		config.GetConfig().Mysql.Host,
		config.GetConfig().Mysql.Database,
	))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.GetConfig().Mysql.TablePrefix + defaultTableName
	}
	Mysql.LogMode(true)
	Mysql.SingularTable(true)
	Mysql.DB().SetMaxIdleConns(10)
	Mysql.DB().SetMaxOpenConns(100)
	Mysql.Callback().Create().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		if !scope.HasError() {
			if createTimeField, ok := scope.FieldByName("CreateTime"); ok {
				if createTimeField.IsBlank {
					createTimeField.Set(time.Now())
				}
			}
			if modifyTimeField, ok := scope.FieldByName("UpdateTime"); ok {
				if modifyTimeField.IsBlank {
					modifyTimeField.Set(time.Now())
				}
			}
		}
	})
	Mysql.Callback().Update().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		if _, ok := scope.Get("gorm:update_column"); !ok {
			scope.SetColumn("UpdateTime", time.Now())
		}
	})
}

func (m *zmInflux) getWriteApi(org, bucket string) (result api.WriteApi) {
	key := org + " " + bucket
	m.Lock()
	defer m.Unlock()
	var ok bool
	if result, ok = m.WriteApiCache[key]; ok {
		return
	} else {
		result = m.Client.WriteApi(org, bucket)
		m.WriteApiCache[key] = result
	}
	return
}
func (m *zmInflux) getQueryApi(org string) (result api.QueryApi) {
	key := org
	m.Lock()
	defer m.Unlock()
	var ok bool
	if result, ok = m.ReadApiCache[key]; ok {
		return
	} else {
		result = m.Client.QueryApi(org)
		m.ReadApiCache[key] = result
	}
	return
}

func (m *zmInflux) Write(bucket string, metric ...*write.Point) (err error) {
	writeApi := m.getWriteApi(config.GetConfig().InfluxDB.Org, bucket)
	defer writeApi.Flush()
	for _, v := range metric {
		writeApi.WritePoint(v)
	}
	return
}

func (m *zmInflux) QueryToRaw(flux string) (raw []byte, err error) {
	readApi := m.Client.QueryApi(config.GetConfig().InfluxDB.Org)
	data, err := readApi.QueryRaw(context.Background(), flux, influxdb2.DefaultDialect())
	raw = []byte(data)
	return
}

func (m *zmInflux) QueryToArray(flux string) (result []map[string]interface{}, err error) {
	readApi := m.getQueryApi(config.GetConfig().InfluxDB.Org)
	data, err := readApi.Query(context.Background(), flux)
	if err != nil {
		return
	}
	for data.Next() {
		result = append(result, data.Record().Values())
	}
	return
}
