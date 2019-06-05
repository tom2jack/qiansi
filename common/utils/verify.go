package utils

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/mojocn/base64Captcha"
	"math/rand"
	"time"
	"tools-server/common/aliyun"
	"tools-server/models"
)

type VerifyStore struct {
	Pool *redis.Pool
}

func init() {
	// init redis store
	customeStore := VerifyStore{models.ZM_Redis}
	base64Captcha.SetCustomStore(&customeStore)
}

// customizeRdsStore implementing Set method of  Store interface
func (s *VerifyStore) Set(id string, value string) {
	conn := s.Pool.Get()
	defer conn.Close()
	conn.Do("SET", id, value, time.Minute*10)
}

// customizeRdsStore implementing Get method of  Store interface
func (s *VerifyStore) Get(id string, clear bool) string {
	conn := s.Pool.Get()
	defer conn.Close()
	reply, _ := redis.Bytes(conn.Do("GET", id))
	if clear {
		conn.Do("DEL", id)
	}
	return string(reply)
}

func VerifyBySMS(phone string) error {
	rand.Seed(time.Now().Unix())
	rnd := rand.Intn(4)
	result := aliyun.SendSmsVerify(phone, string(rnd))
	if !result {
		return fmt.Errorf("发送失败")
	}
	return nil
}

func VerifyByImg(idkey string) (string, string) {
	idkey, captcaInterfaceInstance := base64Captcha.GenerateCaptcha(idkey, base64Captcha.ConfigDigit{})
	base64blob := base64Captcha.CaptchaWriteToBase64Encoding(captcaInterfaceInstance)
	return idkey, base64blob
}

func VerifyCheck(idkey string, value string) bool {
	return base64Captcha.VerifyCaptchaAndIsClear(idkey, value, true)
}
