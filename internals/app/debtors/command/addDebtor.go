package command

import (
	"github.com/Pr3c10us/gbt/internals/domain/debtor"
)

type AddDebtorRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=64"`
	PhoneNumber string `json:"phoneNumber" binding:"required,min=6,max=32"`
	UserID      string `json:"userID" binding:"required"`
}

type AddDebtorCommand interface {
	Handle(request AddDebtorRequest) error
}

type addDebtorCommand struct {
	repository debtor.Repository
}

func (command addDebtorCommand) Handle(request AddDebtorRequest) error {
	newDebtor := debtor.Debtor{
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		UserID:      request.UserID,
	}
	err := command.repository.AddDebtor(&newDebtor)
	if err != nil {
		return err
	}
	return nil
}

func NewAddDebtorCommand(repository debtor.Repository) AddDebtorCommand {
	return &addDebtorCommand{
		repository,
	}
}
