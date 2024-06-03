package queries

import (
	"errors"
	"github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func Test_authenticateUserQuery_Handle(t *testing.T) {
	pw := "secretMan"
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name        string
		args        args
		want        identity.User
		wantErr     bool
		ExpectedErr error
	}{
		{
			name: "valid credentials",
			args: args{
				username: "mrMan",
				password: pw,
			},
			want: identity.User{
				ID:               uuid.New().String(),
				Username:         "MrMan",
				Password:         string(hashedPw),
				SecurityQuestion: "who strong",
				SecurityAnswer:   "me",
			},
			wantErr:     false,
			ExpectedErr: nil,
		},
		{
			name: "incorrect password",
			args: args{
				username: "mrMan",
				password: "wrongPassword",
			},
			want: identity.User{
				ID:               uuid.New().String(),
				Username:         "MrMan",
				Password:         string(hashedPw),
				SecurityQuestion: "who strong",
				SecurityAnswer:   "me",
			},
			wantErr:     true,
			ExpectedErr: appError.BadRequest(IncorrectPasswordErr),
		},
		{
			name: "invalid credentials",
			args: args{
				username: "mrMan2",
				password: "secretMan",
			},
			want:        identity.User{},
			wantErr:     false,
			ExpectedErr: appError.NotFound(errors.New("user with username MrMan2 does not exit")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(identity.MockRepository)
			mockRepo.On("GetUser", tt.args.username).Return(&tt.want, tt.ExpectedErr)
			query := &authenticateUserQuery{
				repository: mockRepo,
			}
			user, err := query.Handle(tt.args.username, tt.args.password)
			if tt.wantErr {
				assert.Equal(t, tt.ExpectedErr, err)
			} else {
				assert.Equal(t, tt.ExpectedErr, err)
				assert.Equal(t, tt.want, *user)
			}
		})
	}
}
