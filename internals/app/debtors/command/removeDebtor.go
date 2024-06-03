package command

import (
	"github.com/Pr3c10us/gbt/internals/domain/debtor"
)

type RemoveDebtorCommand interface {
	Handle(id string) error
}

type removeDebtorCommand struct {
	repository debtor.Repository
}

func (command removeDebtorCommand) Handle(id string) error {
	err := command.repository.RemoveDebtor(id)
	if err != nil {
		return err
	}
	return nil
}

func NewRemoveDebtorCommand(repository debtor.Repository) RemoveDebtorCommand {
	return &removeDebtorCommand{
		repository,
	}
}
