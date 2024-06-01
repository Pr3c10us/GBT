package identity

import (
	"github.com/Pr3c10us/gbt/internals/app/identity/commands"
	"github.com/Pr3c10us/gbt/internals/app/identity/queries"
	"github.com/Pr3c10us/gbt/internals/domain/identity"
)

type Service struct {
	Commands
	Queries
}

type Commands struct {
	CreateUser commands.CreateUserCommand
}

type Queries struct {
	AuthenticateUser queries.AuthenticateUserQuery
}

func NewIdentityService(repository identity.Repository) Service {
	return Service{
		Commands: Commands{
			CreateUser: commands.NewCreateUserCommand(repository),
		},
		Queries: Queries{
			AuthenticateUser: queries.NewAuthenticateUserQuery(repository),
		},
	}
}
