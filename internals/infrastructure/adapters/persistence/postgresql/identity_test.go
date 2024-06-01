package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/Pr3c10us/gbt/pkg/logger"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var (
	sugarLogger = logger.NewSugarLogger(false)
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		message := fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err)
		sugarLogger.Log("fatal", message)
	}

	return db, mock
}

func Test_postgresqlRepository_CreateUser(t *testing.T) {

	tests := []struct {
		name    string
		user    *identity.User
		wantErr bool
	}{
		// Error Case
		{
			name: "Create user error",
			user: &identity.User{
				Username:         "MrMan",
				Password:         "SafeMan",
				SecurityQuestion: "Who Strong",
				SecurityAnswer:   "Me",
			},
			wantErr: true,
		},
		//	Success Case
		{
			name: "Create user success",
			user: &identity.User{
				Username:         "MrMan",
				Password:         "SafeMan",
				SecurityQuestion: "Who Strong",
				SecurityAnswer:   "Me",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get mock and mock db conn
			db, mock := NewMock()
			// create new repo with mock db conn
			repository := &IdentityRepositoryPG{
				db:     db,
				logger: sugarLogger,
			}
			// Close db connection after completion of testcase
			defer func(repository *IdentityRepositoryPG) {
				err := repository.db.Close()
				if err != nil {
					sugarLogger.Log("fatal", "Error closing database connection")
				}
			}(repository)

			// Begin Transaction in mock
			mock.ExpectBegin()
			query := regexp.QuoteMeta(
				`INSERT INTO users (username, password, security_question, security_answer) VALUES ($1,$2,$3,$4)`)
			prep := mock.ExpectPrepare(query)

			if tt.wantErr {
				prep.ExpectExec().
					WithArgs(tt.user.Username, tt.user.Password, tt.user.SecurityQuestion, tt.user.SecurityAnswer).
					WillReturnError(&pq.Error{
						Code:     "23505",
						Severity: "ERROR",
						Message:  "duplicate key value violates unique constraint \"username\"",
						Detail:   "Key (username)=(MrMan) already exists."})
			} else {
				prep.ExpectExec().
					WithArgs(tt.user.Username, tt.user.Password, tt.user.SecurityQuestion, tt.user.SecurityAnswer).
					WillReturnResult(sqlmock.NewResult(0, 1))
			}

			mock.ExpectCommit()

			err := repository.CreateUser(tt.user)
			if (err != nil) != tt.wantErr {
				assert.Equal(t, err, tt.wantErr)
			}
		})
	}
}

func TestIdentityRepositoryPG_GetUser(t *testing.T) {

	type args struct {
		username string
	}
	tests := []struct {
		name        string
		args        args
		want        identity.User
		wantErr     bool
		expectedErr error
	}{
		{
			name: "get user successfully",
			args: args{
				username: "MrMan",
			},
			want: identity.User{
				ID:               uuid.New().String(),
				Username:         "MrMan",
				Password:         "secretMan",
				SecurityQuestion: "who strong",
				SecurityAnswer:   "me",
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "user not found",
			args: args{
				username: "MrMan",
			},
			want:        identity.User{},
			wantErr:     true,
			expectedErr: appError.NotFound(errors.New("user with username MrMan does not exit")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := NewMock()
			repository := &IdentityRepositoryPG{
				db:     db,
				logger: sugarLogger,
			}
			// Close db connection after completion of testcase
			defer func(repository *IdentityRepositoryPG) {
				err := repository.db.Close()
				if err != nil {
					sugarLogger.Log("fatal", "Error closing database connection")
				}
			}(repository)

			query := regexp.QuoteMeta(`SELECT id,username,password,security_question,security_answer FROM users where username=$1`)
			prep := mock.ExpectPrepare(query)
			if tt.wantErr {
				prep.ExpectQuery().WithArgs(tt.args.username).WillReturnError(tt.expectedErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "security_question", "security_answer"}).
					AddRow(tt.want.ID, tt.want.Username, tt.want.Password, tt.want.SecurityQuestion, tt.want.SecurityAnswer)
				prep.ExpectQuery().WithArgs(tt.args.username).WillReturnRows(rows)
			}

			user, err := repository.GetUser(tt.args.username)
			assert.Equal(t, err, tt.expectedErr)
			assert.Equal(t, tt.want, user)
		})
	}
}
