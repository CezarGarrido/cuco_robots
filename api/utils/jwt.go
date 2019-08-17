package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CezarGarrido/cuco_robots/api/entities"
	"github.com/dgrijalva/jwt-go"
)

func ValidToken(r *http.Request) (*entities.Claims, error) {
	claim := &entities.Claims{}

	tokenHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(tokenHeader) != 2 || tokenHeader[0] != "bearer" && tokenHeader[0] != "Bearer" {
		return claim, errors.New("Bearer token inválido.")
	}
	tokenPart := tokenHeader[1]
	token, err := jwt.ParseWithClaims(tokenPart, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte("aplicativo_uems_dourados"), nil
	})
	if !token.Valid {
		return claim, errors.New("Bearer token inválido.")
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claim, errors.New("Assinatura do token inválida.")
		}
		return claim, errors.New("Bearer token inválido.")
	}
	return claim, nil
}
