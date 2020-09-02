package captcha

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/qiansi/common/sdk"
	"github.com/zhi-miao/qiansi/models"
)

type verifyStore struct {
	prefix string
	// Expiration time of captchas.
	expiration time.Duration
}

// 验证码存储库
var VerifyStore *verifyStore

func init() {
	// init redis store
	VerifyStore = &verifyStore{
		"QIANSI:verify:",
		30 * 60,
	}
}

// customizeRdsStore implementing Set method of  Store interface
func (s *verifyStore) Set(id string, value string) {
	models.Redis.Set(context.Background(), s.prefix+id, value, s.expiration)
}

// customizeRdsStore implementing Get method of  Store interface
func (s *verifyStore) Get(id string, clear bool) string {
	reply, _ := models.Redis.Get(context.Background(), s.prefix+id).Result()
	if clear {
		models.Redis.Del(context.Background(), s.prefix+id)
	}
	return reply
}

// Verify 校验功能
func (s *verifyStore) Verify(id, answer string, clear bool) (match bool) {
	reply, _ := models.Redis.Get(context.Background(), s.prefix+id).Result()
	match = reply == answer
	if clear {
		models.Redis.Del(context.Background(), s.prefix+id)
	}
	return
}

// 获取手机短信验证码的idkey
func VerifyBySMSIDKEY(phone string) string {
	return "phone:" + phone
}

// 短信验证码发送
func VerifyBySMS(phone string) error {
	idkey := VerifyBySMSIDKEY(phone)
	lock := VerifyStore.Get(idkey, false)
	if lock != "" {
		return fmt.Errorf("短信已发送，请耐心等待")
	}
	rnd := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	result := sdk.NewAliyunSDK().SendSmsVerify(phone, rnd)
	if !result {
		return fmt.Errorf("发送失败")
	}
	logrus.Infof("短信验证码：%s", rnd)
	VerifyStore.Set(idkey, rnd)
	return nil
}

// 图片验证码生成
func VerifyByImg(idkey string) (string, string) {
	c := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, VerifyStore)
	id, b64s, _ := c.Generate()
	return id, b64s
}

// 验证码校验
func VerifyCheck(idkey string, value string) bool {
	return VerifyStore.Verify(idkey, value, true)
}
