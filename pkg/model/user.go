package model

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/wilsonth122/money-tracker-api/pkg/config"
)

// Token JWT Claims struct
type Token struct {
	UserID string
	jwt.StandardClaims
}

// User struct
type User struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
	Token    string `bson:"token" json:"token"`
}

// GenerateToken - Generates and signs a new JWT
func GenerateToken(id string) string {
	conf := config.New()

	tk := &Token{UserID: id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(conf.Auth.TokenPassword))

	return tokenString
}
