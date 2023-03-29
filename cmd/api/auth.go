package main

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Email string
	Name string 
}

func GenerateJWT(email string, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		Email: email,
		Name: name,
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func Validate(token string) (JWTClaims, error) {
	var jwtClaims JWTClaims
	parsedToken, err := jwt.ParseWithClaims(token, &jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !parsedToken.Valid {
		return jwtClaims, err
	}

	 if !parsedToken.Valid {
		return jwtClaims, errors.New("Invalid token")
	 }

	 return jwtClaims, nil
}
