package models

import (
	"encoding/json"
	"fmt"
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

	ZM_Mysql.SingularTable(true)
	ZM_Mysql.DB().SetMaxIdleConns(10)
	ZM_Mysql.DB().SetMaxOpenConns(100)

	ZM_Mysql.Callback().Create().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		if !scope.HasError() {
			nowTime := time.Now().Unix()
			if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
				if createTimeField.IsBlank {
					createTimeField.Set(nowTime)
				}
			}

			if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
				if modifyTimeField.IsBlank {
					modifyTimeField.Set(nowTime)
				}
			}
		}
	})
	ZM_Mysql.Callback().Update().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {
		if _, ok := scope.Get("gorm:update_column"); !ok {
			scope.SetColumn("ModifiedOn", time.Now().Unix())
		}
	})
	ZM_Mysql.Callback().Delete().Replace("gorm:delete", func(scope *gorm.Scope) {
		addExtraSpaceIfExist := func(str string) string {
			if str != "" {
				return " " + str
			}
			return ""
		}
		if !scope.HasError() {
			var extraOption string
			if str, ok := scope.Get("gorm:delete_option"); ok {
				extraOption = fmt.Sprint(str)
			}

			deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

			if !scope.Search.Unscoped && hasDeletedOnField {
				scope.Raw(fmt.Sprintf(
					"UPDATE %v SET %v=%v%v%v",
					scope.QuotedTableName(),
					scope.Quote(deletedOnField.DBName),
					scope.AddToVars(time.Now().Unix()),
					addExtraSpaceIfExist(scope.CombinedConditionSql()),
					addExtraSpaceIfExist(extraOption),
				)).Exec()
			} else {
				scope.Raw(fmt.Sprintf(
					"DELETE FROM %v%v%v",
					scope.QuotedTableName(),
					addExtraSpaceIfExist(scope.CombinedConditionSql()),
					addExtraSpaceIfExist(extraOption),
				)).Exec()
			}
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
