package queries

import (
	"errors"
	"github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/Pr3c10us/gbt/pkg/utils"
)

var (
	IncorrectPasswordErr = errors.New("incorrect password")
)

type AuthenticateUserQuery interface {
	Handle(username string, password string) (*identity.User, error)
}

type authenticateUserQuery struct {
	repository identity.Repository
}

func NewAuthenticateUserQuery(repository identity.Repository) AuthenticateUserQuery {
	return &authenticateUserQuery{
		repository: repository,
	}
}

func (query *authenticateUserQuery) Handle(username string, password string) (*identity.User, error) {
	user, err := query.repository.GetUser(username)
	if err != nil {
		return &identity.User{}, err
	}

	//	check if password is correct
	if !utils.IsValidPassword(user.Password, password) {
		return &identity.User{}, appError.BadRequest(IncorrectPasswordErr)
	}

	return user, nil
}
