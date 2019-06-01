package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
	"tools-server/conf"
)

// MD5 md5 encryption
func MD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

//ParseToken 解析jwtToken
func ParseToken(tokenString string) (string, error) {
	jwt_secret, err := conf.App.GetValue("app", "jwt_secret")
	if err != nil {
		return "", err
	}
	jwtSecret := []byte(jwt_secret)
	tokenClaims, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
	)

	if tokenClaims == nil {
		if claims, ok := tokenClaims.Claims.(*jwt.StandardClaims); ok && tokenClaims.Valid {
			if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
				return "", fmt.Errorf("过期了")
			}
			return claims.Subject, nil
		}
	}

	return "", fmt.Errorf("解析失败")
}

//CreateToken 生成jwtToken
func CreateToken(subject string, expire time.Duration) (string, error) {
	jwt_secret, err := conf.App.GetValue("app", "jwt_secret")
	if err != nil {
		return "", err
	}
	jwtSecret := []byte(jwt_secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   subject,
		ExpiresAt: time.Now().Add(expire).Unix(),
		Issuer:    "zhimiao-tools-server",
	})
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}
