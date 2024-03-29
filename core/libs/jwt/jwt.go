package jwt

import (
	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
	secretKey []byte
}

func NewJwt(secret string) *Jwt {
	return &Jwt{
		secretKey: []byte(secret),
	}
}

func (this *Jwt) Sign(claims jwt.MapClaims) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	tokenString, err := token.SignedString(this.secretKey)
	if err != nil {
		logger.Error("jwt.Sign", zap.Error(err))
		return ""
	}
	return tokenString
}

func (this *Jwt) Parse(tokenString string) jwt.MapClaims {
	token, err := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return this.secretKey, nil
	})

	if !token.Valid {
		logger.Error("jwt.Parse token not valid")
		return nil
	}

	if err != nil {
		logger.Error("jwt.Parse", zap.Error(err))
		return nil
	}
	return token.Claims.(jwt.MapClaims)
}
