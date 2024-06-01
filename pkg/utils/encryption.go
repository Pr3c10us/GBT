package utils

import "golang.org/x/crypto/bcrypt"

func IsValidPassword(hashedPw string, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw)) == nil
}
