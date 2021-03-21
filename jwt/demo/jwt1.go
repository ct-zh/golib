package demo

import (
	"github.com/dgrijalva/jwt-go"
)

// 根据uid与secret创建token
// 这里加密方法使用HS256
func CreateToken(uid, secret string, timeout int64) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": timeout,
	})

	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// 根据token与secret解析token
func ParseToken(token, secret string) (string, error) {
	cliam, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	return cliam.Claims.(jwt.MapClaims)["uid"].(string), nil
}
