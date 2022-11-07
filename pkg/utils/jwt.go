package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SecretKey 认证key
const (
	SecretKey = "gintemplate"
	Issuer    = "gintemplate"
)

// jwtCustomClaims token签名信息
type jwtCustomClaims struct {
	jwt.StandardClaims

	Uid      string
	ClientIp string
}

// GenerateToken 生成token
func GenerateToken(uid, clientIp string) (string, error) {
	claims := &jwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			Issuer:    Issuer,
		},
		Uid:      uid,
		ClientIp: clientIp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyToken 验证token
func VerifyToken(tokenSrt string) (*jwtCustomClaims, error) {

	var claims jwtCustomClaims

	verifyToken, err := jwt.ParseWithClaims(tokenSrt, &claims, func(*jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !verifyToken.Valid {
		return nil, err
	}

	return &claims, nil
}
