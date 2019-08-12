package entities

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Aluno `json:"aluno"`
	jwt.StandardClaims
}
