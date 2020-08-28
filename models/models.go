package models

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	"github.com/influxdata/influxdb-client-go/api/write"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zhi-miao/qiansi/common"
)

var (
	Redis    *zmRedis
	Mysql    *gorm.DB
	InfluxDB *zmInflux
)

type zmRedis struct {
	redis.Pool
}

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
		Client:        influxdb2.NewClient(common.Config.InfluxDB.Host, common.Config.InfluxDB.Token),
		WriteApiCache: make(map[string]api.WriteApi),
		ReadApiCache:  make(map[string]api.QueryApi),
	}
}

// Setup Initialize the Redis instance
func loadRedis() {
	Redis = &zmRedis{
		redis.Pool{
			MaxIdle:     common.Config.Redis.MaxIdle,
			MaxActive:   common.Config.Redis.MaxActive,
			IdleTimeout: time.Duration(common.Config.Redis.IdleTimeOut) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", common.Config.Redis.Host)
				if err != nil {
					return nil, err
				}
				if common.Config.Redis.Auth != "" {
					if _, err := c.Do("AUTH", common.Config.Redis.Auth); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

// Setup Initialize the Mysql instance
func loadMysql() {
	var err error
	Mysql, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		common.Config.Mysql.User,
		common.Config.Mysql.Password,
		common.Config.Mysql.Host,
		common.Config.Mysql.Database,
	))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return common.Config.Mysql.TablePrefix + defaultTableName
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

// Set a key/value
func (r *zmRedis) Set(key string, value string, time int) error {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

// Exists check a key
func (r *zmRedis) Exists(key string) bool {
	conn := r.Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func (r *zmRedis) Get(key string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return reply, nil
}

// Delete delete a kye
func (r *zmRedis) Del(key string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("DEL", key))
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
	writeApi := m.getWriteApi(common.Config.InfluxDB.Org, bucket)
	defer writeApi.Flush()
	for _, v := range metric {
		writeApi.WritePoint(v)
	}
	return
}

func (m *zmInflux) QueryToRaw(flux string) (raw []byte, err error) {
	readApi := m.Client.QueryApi(common.Config.InfluxDB.Org)
	data, err := readApi.QueryRaw(context.Background(), flux, influxdb2.DefaultDialect())
	raw = []byte(data)
	return
}

func (m *zmInflux) QueryToArray(flux string) (result []map[string]interface{}, err error) {
	readApi := m.getQueryApi(common.Config.InfluxDB.Org)
	data, err := readApi.Query(context.Background(), flux)
	if err != nil {
		return
	}
	for data.Next() {
		result = append(result, data.Record().Values())
	}
	return
}
