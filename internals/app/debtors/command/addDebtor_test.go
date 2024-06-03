package command

import (
	"github.com/Pr3c10us/gbt/internals/domain/debtor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_addDebtorCommand_Handle(t *testing.T) {
	type args struct {
		request AddDebtorRequest
	}
	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "successful addition of debtor",
			args: args{
				request: AddDebtorRequest{
					Name:        "mrWoman",
					PhoneNumber: "07000000000",
					UserID:      "0000000000",
				},
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(debtor.MockRepository)
			mockDebtor := debtor.Debtor{
				Name:        tt.args.request.Name,
				PhoneNumber: tt.args.request.PhoneNumber,
				UserID:      tt.args.request.UserID,
			}
			mockRepo.On("AddDebtor", &mockDebtor).Return(tt.expectedErr)

			command := &addDebtorCommand{
				repository: mockRepo,
			}
			assert.Equal(t, tt.expectedErr, command.Handle(tt.args.request))
		})
	}
}
