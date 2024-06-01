package commands

import (
	"github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/Pr3c10us/gbt/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateUserCommand_Handle(t *testing.T) {
	type args struct {
		addUserRequest identity.AddUserRequest
	}
	tests := []struct {
		name string
		args args
		Err  error
	}{
		{
			name: "valid input",
			args: args{
				addUserRequest: identity.AddUserRequest{
					Username:         "MrMan",
					Password:         "secretMan",
					SecurityQuestion: "who strong",
					SecurityAnswer:   "me",
				},
			},
			Err: nil,
		},
		{
			name: "invalid input",
			args: args{
				addUserRequest: identity.AddUserRequest{
					Username:         "MrMan",
					Password:         "",
					SecurityQuestion: "who strong",
					SecurityAnswer:   "me",
				},
			},
			Err: appError.BadRequest(identity.InvalidPasswordError),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			mockRepo := new(identity.MockRepository)
			user, _ := identity.NewUser(tt.args.addUserRequest)
			mockRepo.On("CreateUser", mock.MatchedBy(utils.UserMatcher(user))).Return(tt.Err)
			command := &createUserCommand{
				repository: mockRepo,
			}
			a.Equal(tt.Err, command.Handle(tt.args.addUserRequest), "the errors should be the same")
		})
	}
}
