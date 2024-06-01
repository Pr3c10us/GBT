package utils

import (
	"errors"
	"github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateTokenFromUser(user *identity.User, secret string) (string, error) {
	now := time.Now()
	expires := now.Add(time.Hour * 24 * 30).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", appError.InternalServerError(errors.New("failed to authenticate user"))
	}
	return tokenStr, nil
}
