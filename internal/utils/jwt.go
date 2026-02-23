package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("super-secret-key") //config

type Claims struct {
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"` 
	jwt.RegisteredClaims
}

/*
====================================
GENERATE TOKEN
====================================
*/

func GenerateToken(userID int64, roles []string) (string, error) {

	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtSecret)
}

/*
====================================
PARSE TOKEN
====================================
*/

func ParseToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}