package tools

import (
	"github.com/golang-jwt/jwt/v5"
	"tiktok/global"
	"time"
)

func GenerateToken(username string) (string, error) {
	strings := make([]string, 1)
	strings[0] = username
	jwtConfig := global.Config.JWT
	claims := jwt.RegisteredClaims{
		Issuer:    jwtConfig.Issuer,                                                                // 发行人 、 签发者
		Subject:   jwtConfig.Subject,                                                               // 主题 、 面向的用户
		Audience:  strings,                                                                         // 用户、接收者
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConfig.TokenExpireDuration * time.Second)), // 到期时间  必须大于签发时间
		IssuedAt:  jwt.NewNumericDate(time.Now()),                                                  // 发布时间，签发时间
		ID:        "",                                                                              // jwt的唯一身份标识
		NotBefore: nil,                                                                             // 在此之前不可用
	}
	signedString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtConfig.Secrete))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func ParseToken(token string) (*jwt.RegisteredClaims, error) {
	claims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.Secrete), nil
	})
	registeredClaims := claims.Claims.(*jwt.RegisteredClaims)
	if claims.Valid {
		return registeredClaims, nil
	}
	return nil, err
}
