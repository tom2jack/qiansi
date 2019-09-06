package models

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"qiansi/common/conf"
	"time"
)

var (
	ZM_Redis *Redis
	ZM_Mysql *gorm.DB
)

type Redis struct {
	redis.Pool
}

type CommonMap map[string]interface{}

type PageParam struct {
	LastId int
	Page int `binding:"min=1"`
	PageSize int `binding:"min=1,max=50"`
}
type PageInfo struct {
	Page int
	PageSize int
	TotalSize int
	Rows interface{}
}

type ModelBase1 struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	UpdateTime time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	CreateTime time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}

func init()  {
	loadRedis()
	loadMysql()
}

func (p *PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Setup Initialize the Redis instance
func loadRedis() {
	ZM_Redis = &Redis{
		redis.Pool{
			MaxIdle:     conf.S.MustInt("redis", "max_idle"),
			MaxActive:   conf.S.MustInt("redis", "max_active"),
			IdleTimeout: time.Duration(conf.S.MustInt("redis", "idle_timeout")) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", conf.S.MustValue("redis", "host"))
				if err != nil {
					return nil, err
				}
				if conf.S.MustValue("redis", "auth") != "" {
					if _, err := c.Do("AUTH", conf.S.MustValue("redis", "auth")); err != nil {
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
	ZM_Mysql, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.S.MustValue("mysql", "user"),
		conf.S.MustValue("mysql", "password"),
		conf.S.MustValue("mysql", "host"),
		conf.S.MustValue("mysql", "database"),
	))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.S.MustValue("mysql", "table_prefix") + defaultTableName
	}

	ZM_Mysql.LogMode(true)
	ZM_Mysql.SingularTable(true)
	ZM_Mysql.DB().SetMaxIdleConns(10)
	ZM_Mysql.DB().SetMaxOpenConns(100)

	ZM_Mysql.Callback().Create().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
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
	ZM_Mysql.Callback().Update().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		if _, ok := scope.Get("gorm:update_column"); !ok {
			scope.SetColumn("UpdateTime", time.Now())
		}
	})
}

// Set a key/value
func (r *Redis) Set(key string, value string, time int) error {
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
func (r *Redis) Exists(key string) bool {
	conn := r.Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func (r *Redis) Get(key string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return reply, nil
}

// Delete delete a kye
func (r *Redis) Del(key string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("DEL", key))
}
