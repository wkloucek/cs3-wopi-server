package app

import "github.com/golang-jwt/jwt"

type Claims struct {
	WopiContext WopiContext `json:"WopiContext"`
	jwt.StandardClaims
}
