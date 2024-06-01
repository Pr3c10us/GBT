package commands

import (
	"github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/pkg/appError"
)

type CreateUserCommand interface {
	Handle(request identity.AddUserRequest) error
}

type createUserCommand struct {
	repository identity.Repository
}

func NewCreateUserCommand(repository identity.Repository) CreateUserCommand {
	return &createUserCommand{
		repository: repository,
	}
}

func (command *createUserCommand) Handle(request identity.AddUserRequest) error {
	user, err := identity.NewUser(request)
	if err != nil {
		return appError.BadRequest(err)
	}
	err = command.repository.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
