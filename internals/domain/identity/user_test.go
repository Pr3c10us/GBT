package identity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		addUserRequest AddUserRequest
	}
	tests := []struct {
		name string
		args args
		Err  error
	}{
		{
			name: "invalid username",
			args: args{
				AddUserRequest{
					Username:         "",
					Password:         "secretMan",
					SecurityQuestion: "who strong",
					SecurityAnswer:   "Me",
				},
			},
			Err: InvalidUsernameError,
		},
		{
			name: "invalid password",
			args: args{
				AddUserRequest{
					Username:         "mrMan",
					Password:         "",
					SecurityQuestion: "who strong",
					SecurityAnswer:   "Me",
				},
			},
			Err: InvalidPasswordError,
		},
		{
			name: "invalid security question",
			args: args{
				AddUserRequest{
					Username:         "mrMan",
					Password:         "secretMan",
					SecurityQuestion: "",
					SecurityAnswer:   "Me",
				},
			},
			Err: InvalidQuestionError,
		},
		{
			name: "invalid security answer",
			args: args{
				AddUserRequest{
					Username:         "mrMan",
					Password:         "secretMan",
					SecurityQuestion: "who strong",
					SecurityAnswer:   "",
				},
			},
			Err: InvalidAnswerError,
		},
		{
			name: "valid params",
			args: args{
				AddUserRequest{
					Username:         "mrMan",
					Password:         "secretMan",
					SecurityQuestion: "who strong",
					SecurityAnswer:   "Me",
				},
			},
			Err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			_, err := NewUser(tt.args.addUserRequest)
			a.Equal(err, tt.Err, "the errors should be the same")

		})
	}
}
