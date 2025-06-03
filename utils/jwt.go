package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

func GenerateToken(email string) (string, error) {

	// with jwt claims simply mean the data that's attached to it.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	// key should be bytes slice
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (string, error) {

	// keyFn determine how to verify the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// check signing method is HMAC
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method.")
		}

		// return verified signing key
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", errors.New("Could not parse token.")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return "", errors.New("Invalid token.")
	}

	// extract claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("Invalid token claims.")
	}

	// extract email info
	userEmail := claims["email"].(string)

	return userEmail, nil
}
