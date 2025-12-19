package jwt

import (
	"os"
)

var JWTSecret string

func LoadJWT() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "jaimaheshbabaubabulakebabusreemanthudu"
	}
	return secret
}
