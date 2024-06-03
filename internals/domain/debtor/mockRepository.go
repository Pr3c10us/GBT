package debtor

import "github.com/stretchr/testify/mock"

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetDebtorByID(id string) (*Debtor, error) {
	args := m.Called(id)
	return args.Get(0).(*Debtor), args.Error(1)
}

func (m *MockRepository) AddDebtor(debtor *Debtor) error {
	args := m.Called(debtor)
	return args.Error(0)
}

func (m *MockRepository) RemoveDebtor(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
