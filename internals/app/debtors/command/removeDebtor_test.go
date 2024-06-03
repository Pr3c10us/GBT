package command

import (
	"github.com/Pr3c10us/gbt/internals/domain/debtor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_removeDebtorCommand_Handle(t *testing.T) {

	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "successful removal of debtor",
			args: args{
				id: "111111111",
			},
			expectedErr: nil,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(debtor.MockRepository)
			mockRepo.On("RemoveDebtor", tt.args.id).Return(tt.expectedErr)

			command := &removeDebtorCommand{
				repository: mockRepo,
			}
			assert.Equal(t, tt.expectedErr, command.Handle(tt.args.id))
		})
	}
}
