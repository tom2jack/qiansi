package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.StandardClaims); ok && tokenClaims.Valid {
			if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
				return "", fmt.Errorf("过期了")
			}
			if claims.Issuer != "zhimiao-tools-server" {
				return "", fmt.Errorf("非法来源的签名")
			}
			return claims.Subject, nil
		}
	}
	return "", err
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

// 密码生成
func PasswordHash(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// 密码验证
func PasswordVerify(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}
