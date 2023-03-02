package tools

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string) (string, error) {
	strings := make([]string, 1)
	strings[0] = username
	claims := jwt.RegisteredClaims{
		Issuer:    "",      // 发行人 、 签发者
		Subject:   "",      // 主题 、 面向的用户
		Audience:  strings, // 用户、接收者
		ExpiresAt: nil,     // 到期时间  必须大于签发时间
		NotBefore: nil,     // 在此之前不可用
		IssuedAt:  nil,     // 发布时间，签发时间
		ID:        "",      // jwt的唯一身份标识
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("AllYourBase"))
}

func ParseToken(token string) (*jwt.RegisteredClaims, error) {
	claims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})
	registeredClaims := claims.Claims.(*jwt.RegisteredClaims)
	if claims.Valid {
		return registeredClaims, nil
	}
	return nil, err
}
