package identity

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID               string
	Username         string
	Password         string
	SecurityQuestion string
	SecurityAnswer   string
}

type AddUserRequest struct {
	Username         string `json:"username" binding:"required,max=16"`
	Password         string `json:"password" binding:"required,max=32"`
	SecurityQuestion string `json:"securityQuestion" binding:"required"`
	SecurityAnswer   string `json:"securityAnswer" binding:"required"`
}

var (
	InvalidUsernameError = errors.New("username provided is invalid")
	InvalidPasswordError = errors.New("password provided is invalid")
	InvalidQuestionError = errors.New("security question provided is invalid")
	InvalidAnswerError   = errors.New("security answer provided is invalid")
)

func NewUser(addUserRequest AddUserRequest) (*User, error) {
	if len(addUserRequest.Username) < 1 || len(addUserRequest.Username) > 16 {
		return &User{}, InvalidUsernameError
	}
	if len(addUserRequest.Password) < 1 || len(addUserRequest.Password) > 36 {
		return &User{}, InvalidPasswordError
	}
	if len(addUserRequest.SecurityQuestion) < 1 || len(addUserRequest.SecurityQuestion) > 255 {
		return &User{}, InvalidQuestionError
	}
	if len(addUserRequest.SecurityAnswer) < 1 || len(addUserRequest.SecurityAnswer) > 255 {
		return &User{}, InvalidAnswerError
	}

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(addUserRequest.Password), 14)
	if err != nil {
		return &User{}, err
	}
	securityAnswerBytes, err := bcrypt.GenerateFromPassword([]byte(strings.ToLower(addUserRequest.SecurityAnswer)), 14)
	if err != nil {
		return &User{}, err
	}

	return &User{
		Username:         addUserRequest.Username,
		Password:         string(passwordBytes),
		SecurityQuestion: strings.ToLower(addUserRequest.SecurityQuestion),
		SecurityAnswer:   string(securityAnswerBytes),
	}, nil
}
