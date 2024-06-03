package debtors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pr3c10us/gbt/internals/app/debtors"
	"github.com/Pr3c10us/gbt/internals/app/debtors/command"
	"github.com/Pr3c10us/gbt/internals/domain/debtor"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/Pr3c10us/gbt/pkg/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_AddDebtor(t *testing.T) {
	type args struct {
		debtor command.AddDebtorRequest
	}
	tests := []struct {
		name        string
		args        args
		statusCode  int
		ExpectedErr error
	}{
		{
			name: "successful request",
			args: args{
				debtor: command.AddDebtorRequest{
					Name:        "mrWoman",
					PhoneNumber: "07000000000",
					UserID:      "445414e9-e4c7-4255-8d02-b1fdf1fbd0f9",
				},
			},
			statusCode:  http.StatusOK,
			ExpectedErr: nil,
		},
		{
			name: "invalid request arguments",
			args: args{
				debtor: command.AddDebtorRequest{
					Name:        "",
					PhoneNumber: "",
					UserID:      "",
				},
			},
			statusCode:  http.StatusNotAcceptable,
			ExpectedErr: nil,
		},
		{
			name: "debtor name taken",
			args: args{
				debtor: command.AddDebtorRequest{
					Name:        "mrWoman",
					PhoneNumber: "07000000000",
					UserID:      "445414e9-e4c7-4255-8d02-b1fdf1fbd0f9",
				},
			},
			statusCode:  http.StatusConflict,
			ExpectedErr: appError.NewPQUniqueError(),
		},
		{
			name: "invalid uuid",
			args: args{
				debtor: command.AddDebtorRequest{
					Name:        "mrWoman",
					PhoneNumber: "07000000000",
					UserID:      "0000000000",
				},
			},
			statusCode:  http.StatusBadRequest,
			ExpectedErr: appError.NewPQSyntaxError(),
		},
		{
			name: "invalid user Id",
			args: args{
				debtor: command.AddDebtorRequest{
					Name:        "mrWoman",
					PhoneNumber: "07000000000",
					UserID:      "445414e9-e4c7-4255-8d02-b1fdf1fbd0f9",
				},
			},
			statusCode:  http.StatusBadRequest,
			ExpectedErr: appError.NewPQForeignError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDebtor := debtor.Debtor{
				Name:        tt.args.debtor.Name,
				PhoneNumber: tt.args.debtor.PhoneNumber,
				UserID:      tt.args.debtor.UserID,
			}

			mockRepo := new(debtor.MockRepository)
			mockRepo.On("AddDebtor", &mockDebtor).Return(tt.ExpectedErr)

			services := debtors.NewDebtorServices(mockRepo)
			handle := NewDebtorsHandler(services)

			engine := utils.SetupRouter()
			engine.POST("/api/v1/debtors", handle.AddDebtor)

			w := httptest.NewRecorder()
			addDebtorJSON, _ := json.Marshal(tt.args.debtor)
			r, _ := http.NewRequest("POST", "/api/v1/debtors", strings.NewReader(string(addDebtorJSON)))
			r.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(w, r)

			assert.Equal(t, tt.statusCode, w.Code)

		})
	}
}

func TestHandler_RemoveDebtor(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		statusCode  int
		ExpectedErr error
	}{
		{
			name: "successful request",
			args: args{
				id: "b79630dc-5283-4260-b033-9193cd83e23b",
			},
			statusCode:  http.StatusOK,
			ExpectedErr: nil,
		},
		{
			name: "debtor not found",
			args: args{
				id: "b79630dc-5283-4260-b033-9193cd83e23a",
			},
			statusCode:  http.StatusBadRequest,
			ExpectedErr: appError.BadRequest(errors.New("no debtor deleted: the specified debtor id does not exist")),
		},
		{
			name: "invalid uuid",
			args: args{
				id: "0000000000",
			},
			statusCode:  http.StatusNotAcceptable,
			ExpectedErr: appError.NewPQSyntaxError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(debtor.MockRepository)
			mockRepo.On("RemoveDebtor", tt.args.id).Return(tt.ExpectedErr)

			services := debtors.NewDebtorServices(mockRepo)
			handle := NewDebtorsHandler(services)

			engine := utils.SetupRouter()
			engine.DELETE("/api/v1/debtors/:id", handle.RemoveDebtor)

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/debtors/%v", tt.args.id), nil)
			engine.ServeHTTP(w, r)

			assert.Equal(t, tt.statusCode, w.Code)

		})
	}
}
