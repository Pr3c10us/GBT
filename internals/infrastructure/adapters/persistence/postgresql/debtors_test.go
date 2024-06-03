package postgresql

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Pr3c10us/gbt/internals/domain/debtor"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/Pr3c10us/gbt/pkg/utils"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestDebtorRepositoryPG_AddDebtor(t *testing.T) {

	type args struct {
		debtor *debtor.Debtor
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid debtor",
			args: args{
				debtor: &debtor.Debtor{
					Name:        "mrWoman",
					PhoneNumber: "070000000000",
					UserID:      "0000000000",
				},
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "invalid user id",
			args: args{
				debtor: &debtor.Debtor{
					Name:        "mrWoman",
					PhoneNumber: "070000000000",
					UserID:      "0000000000",
				},
			},
			wantErr:     true,
			expectedErr: appError.NewPQForeignError(),
		},
		{
			name: "debtor already exist",
			args: args{
				debtor: &debtor.Debtor{
					Name:        "mrWoman",
					PhoneNumber: "070000000000",
					UserID:      "0000000000",
				},
			},
			wantErr:     true,
			expectedErr: appError.NewPQUniqueError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := utils.NewMock()
			repository := &DebtorRepositoryPG{
				db,
				sugarLogger,
			}
			defer func(repository *DebtorRepositoryPG) {
				err := repository.db.Close()
				if err != nil {
					sugarLogger.Log("fatal", "Error closing database connection")
				}
			}(repository)

			query := regexp.QuoteMeta(
				`INSERT INTO debtors ( name, phone_number, user_id) VALUES ($1,$2,$3)`)
			prep := mock.ExpectPrepare(query)

			if tt.wantErr {
				prep.ExpectExec().WithArgs(tt.args.debtor.Name, tt.args.debtor.PhoneNumber, tt.args.debtor.UserID).WillReturnError(tt.expectedErr)
			} else {
				prep.ExpectExec().WithArgs(tt.args.debtor.Name, tt.args.debtor.PhoneNumber, tt.args.debtor.UserID).WillReturnResult(sqlmock.NewResult(0, 1))

			}

			err := repository.AddDebtor(tt.args.debtor)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestDebtorRepositoryPG_GetDebtorByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		want        *debtor.Debtor
		wantErr     bool
		expectedErr error
	}{
		{
			name: "get user successful",
			args: args{
				id: "1111111111",
			},
			want: &debtor.Debtor{
				ID:          "1111111111",
				Name:        "mrWoman",
				PhoneNumber: "07000000000",
				UserID:      "0000000000",
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "user does not exist",
			args: args{
				id: "1111111111",
			},
			want:        &debtor.Debtor{},
			wantErr:     true,
			expectedErr: appError.NotFound(errors.New(fmt.Sprintf("debtor with id '%v' does not exit", "1111111111"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := utils.NewMock()
			repository := &DebtorRepositoryPG{
				db,
				sugarLogger,
			}
			defer func(repository *DebtorRepositoryPG) {
				err := repository.db.Close()
				if err != nil {
					sugarLogger.Log("fatal", "Error closing database connection")
				}
			}(repository)

			query := regexp.QuoteMeta(`SELECT id, name, phone_number, user_id FROM debtors WHERE id=$1;`)
			if tt.wantErr {
				mock.ExpectQuery(query).WithArgs(tt.args.id).WillReturnError(tt.expectedErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name", "phone_number", "user_id"}).AddRow(
					tt.want.ID, tt.want.Name, tt.want.PhoneNumber, tt.want.UserID,
				)
				mock.ExpectQuery(query).WithArgs(tt.args.id).WillReturnRows(rows)
			}

			d, err := repository.GetDebtorByID(tt.args.id)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.want, d)
		})
	}
}

func TestDebtorRepositoryPG_RemoveDebtor(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		affectedRow int64
		expectedErr error
	}{
		{
			name: "debtor deleted succesfully",
			args: args{
				id: "1111111111",
			},
			affectedRow: 1,
			expectedErr: nil,
		},
		{
			name: "debtor deleted succesfully",
			args: args{
				id: "2222222222",
			},
			affectedRow: 0,
			expectedErr: appError.BadRequest(errors.New("no debtor deleted: the specified debtor id does not exist")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := utils.NewMock()
			repository := &DebtorRepositoryPG{
				db,
				sugarLogger,
			}
			defer func(repository *DebtorRepositoryPG) {
				err := repository.db.Close()
				if err != nil {
					sugarLogger.Log("fatal", "Error closing database connection")
				}
			}(repository)

			query := regexp.QuoteMeta(`DELETE FROM debtors WHERE id = $1`)
			prep := mock.ExpectPrepare(query)

			prep.ExpectExec().WithArgs(tt.args.id).WillReturnResult(sqlmock.NewResult(0, tt.affectedRow))

			err := repository.RemoveDebtor(tt.args.id)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
