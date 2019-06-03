package models

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
	"tools-server/conf"
)

var (
	ZM_Redis *redis.Pool
	ZM_Mysql *gorm.DB
)

type ModelBase1 struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreateTime int `json:"created_time"`
	UpdateTime int `json:"update_time"`
}

type Json struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func Show(C *gin.Context, code int, msg string, data interface{}) {
	C.JSON(200, Json{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	return
}

// Setup Initialize the Redis instance
func LoadRedis() {
	ZM_Redis = &redis.Pool{
		MaxIdle:     conf.App.MustInt("redis", "max_idle"),
		MaxActive:   conf.App.MustInt("redis", "max_active"),
		IdleTimeout: time.Duration(conf.App.MustInt("redis", "idle_timeout")) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.App.MustValue("redis", "host"))
			if err != nil {
				return nil, err
			}
			if conf.App.MustValue("redis", "auth") != "" {
				if _, err := c.Do("AUTH", conf.App.MustValue("redis", "auth")); err != nil {
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
	}
}

// Setup Initialize the Mysql instance
func LoadMysql() {
	var err error
	ZM_Mysql, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.App.MustValue("mysql", "user"),
		conf.App.MustValue("mysql", "password"),
		conf.App.MustValue("mysql", "host"),
		conf.App.MustValue("mysql", "database"),
	))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.App.MustValue("mysql", "table_prefix") + defaultTableName
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
func RedisSet(key string, data interface{}, time int) error {
	conn := ZM_Redis.Get()
	defer conn.Close()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, value)
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
func RedisExists(key string) bool {
	conn := ZM_Redis.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func RedisGet(key string) ([]byte, error) {
	conn := ZM_Redis.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func RedisDelete(key string) (bool, error) {
	conn := ZM_Redis.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}
