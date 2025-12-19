package jwt

import (
	// Go Internal Packages
	"fmt"
	"time"

	// External Packages

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	SecretKey     string
	SecretKeyByte []byte
}

func NewJwt(secret string) *Jwt {
	return &Jwt{
		SecretKey:     secret,
		SecretKeyByte: []byte(secret),
	}
}

func (j *Jwt) Decode(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.SecretKeyByte, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidId
	}

	return claims, nil
}

func (j *Jwt) GenerateJwtToken(userId string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user": userId,
		"exp":  time.Now().Add(duration).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)

	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func FetchClaim(key string, claims jwt.MapClaims) string {
	if value, exists := claims[key].(string); exists && value != "" {
		return value
	}
	return ""
}
