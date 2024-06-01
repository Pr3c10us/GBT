package identity

import "github.com/stretchr/testify/mock"

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUser(username string) (User, error) {
	args := m.Called(username)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockRepository) CreateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}
