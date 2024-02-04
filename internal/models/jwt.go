package models

import (
	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	jwt.StandardClaims
	UserId  int  `json:"userId"`
	IsAdmin bool `json:"isAdmin"`
}

type Role int

const (
	Client Role = iota // 0
	Admin              // 1
)

type UserLogin struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserSignUp struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Email string `json:"email"`
}