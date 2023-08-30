package util

import (
	"os"

	"github.com/golang-jwt/jwt"
)

func GenerateJwtToken(issuer string) (string, error) {
	claim := &jwt.StandardClaims{
		Issuer: issuer,
	}

	// NOTE: Encrypt the claim using your private key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	secret := os.Getenv("JWT_SECRET_KEY")

	return token.SignedString([]byte(secret))
}

func ParseJwt(claim string) error {
	// NOTE: Decrypt the cookie value using the private key
	_, err := jwt.ParseWithClaims(claim, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return err
	}
	return nil
}
