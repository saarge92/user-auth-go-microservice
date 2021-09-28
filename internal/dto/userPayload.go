package dto

import "github.com/dgrijalva/jwt-go"

type UserPayLoad struct {
	jwt.StandardClaims
	UserName string   `json:"login"`
	Roles    []string `json:"roles"`
}
