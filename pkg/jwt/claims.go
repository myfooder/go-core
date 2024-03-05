package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId      string `json:"uid"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Locale      string `json:"locale"`
	jwt.RegisteredClaims
}
