package captcha

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"math/rand"
	"qiansi/common/models"
	"qiansi/common/zmlog"
	"time"
)

type VerifyStore struct {
	prefix string
	// Expiration time of captchas.
	expiration int
}

func init() {
	// init redis store
	customeStore := &VerifyStore{
		"QIANSI:verify:",
		30 * 60,
	}
	base64Captcha.SetCustomStore(customeStore)
}

// customizeRdsStore implementing Set method of  Store interface
func (s *VerifyStore) Set(id string, value string) {
	models.ZM_Redis.Set(s.prefix+id, value, s.expiration)
}

// customizeRdsStore implementing Get method of  Store interface
func (s *VerifyStore) Get(id string, clear bool) string {
	reply, _ := models.ZM_Redis.Get(s.prefix + id)
	if clear {
		models.ZM_Redis.Del(s.prefix + id)
	}
	return string(reply)
}

// 获取手机短信验证码的idkey
func VerifyBySMSIDKEY(phone string) string {
	return "QIANSI:verify:phone:" + phone
}

// 短信验证码发送
func VerifyBySMS(phone string) error {
	idkey := VerifyBySMSIDKEY(phone)
	lock := models.ZM_Redis.Exists(idkey)
	if lock {
		return fmt.Errorf("短信已发送，请耐心等待")
	}
	rnd := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	/*result := aliyun.SendSmsVerify(phone, string(rnd))
	if !result {
		return fmt.Errorf("发送失败")
	}*/
	zmlog.Info("短信验证码：%s", rnd)
	models.ZM_Redis.Set(idkey, rnd, 30*60)
	return nil
}

// 图片验证码生成
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
		1,
		// MaxSkew max absolute skew factor of a single digit.
		// 图像验证码的最大干扰洗漱.
		1, //4.5,
		// DotCount Number of background circles.
		// 图像验证码干扰圆点的数量.
		0, //30,
	})
	base64blob := base64Captcha.CaptchaWriteToBase64Encoding(captcaInterfaceInstance)
	return idkey, base64blob
}

// 验证码校验
func VerifyCheck(idkey string, value string) bool {
	return base64Captcha.VerifyCaptchaAndIsClear(idkey, value, false)
}
