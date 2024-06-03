package query

import (
	"github.com/Pr3c10us/gbt/internals/domain/debtor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getDebtorByIDQuery_Handle(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		want        *debtor.Debtor
		expectedErr error
	}{
		{
			name: "successful fetch of debtor",
			args: args{
				id: "1111111111",
			},
			want: &debtor.Debtor{
				ID:          "1111111111",
				Name:        "mrWoman",
				PhoneNumber: "07000000000",
				UserID:      "0000000000",
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(debtor.MockRepository)
			mockRepo.On("GetDebtorByID", tt.args.id).Return(tt.want, tt.expectedErr)

			query := &getDebtorByIDQuery{
				repository: mockRepo,
			}
			d, err := query.Handle(tt.args.id)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.want, d)
			assert.Equal(t, tt.args.id, d.ID)
		})
	}
}
