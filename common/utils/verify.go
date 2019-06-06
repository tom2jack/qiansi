package utils

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"math/rand"
	"time"
	"tools-server/common/aliyun"
	"tools-server/models"
)

type VerifyStore struct {
	prefix string
	// Expiration time of captchas.
	expiration int
}

func init() {
	// init redis store
	customeStore := &VerifyStore{
		"ZMT:verify:",
		30 * 60,
	}
	base64Captcha.SetCustomStore(customeStore)
}

// customizeRdsStore implementing Set method of  Store interface
func (s *VerifyStore) Set(id string, value string) {
	models.RedisSet(s.prefix+id, value, s.expiration)
}

// customizeRdsStore implementing Get method of  Store interface
func (s *VerifyStore) Get(id string, clear bool) string {
	reply, _ := models.RedisGet(s.prefix + id)
	if clear {
		models.RedisDelete(s.prefix + id)
	}
	return string(reply)
}

func VerifyBySMS(phone string) error {
	idkey := "ZMT:verify:phone:" + phone
	lock := models.RedisExists(idkey)
	if lock {
		return fmt.Errorf("短信已发送，请耐心等待")
	}
	rnd := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	result := aliyun.SendSmsVerify(phone, string(rnd))
	if !result {
		return fmt.Errorf("发送失败")
	}
	models.RedisSet(idkey, rnd, 30*60)
	return nil
}

func VerifyByImg(idkey string) (string, string) {
	idkey, captcaInterfaceInstance := base64Captcha.GenerateCaptcha(idkey, base64Captcha.ConfigDigit{
		// Height png height in pixel.
		// 图像验证码的高度像素.
		50,
		// Width Captcha png width in pixel.
		// 图像验证码的宽度像素
		150,
		// DefaultLen Default number of digits in captcha solution.
		// 默认数字验证长度6.
		6,
		// MaxSkew max absolute skew factor of a single digit.
		// 图像验证码的最大干扰洗漱.
		4.5,
		// DotCount Number of background circles.
		// 图像验证码干扰圆点的数量.
		30,
	})
	base64blob := base64Captcha.CaptchaWriteToBase64Encoding(captcaInterfaceInstance)
	return idkey, base64blob
}

func VerifyCheck(idkey string, value string) bool {
	return base64Captcha.VerifyCaptchaAndIsClear(idkey, value, true)
}
