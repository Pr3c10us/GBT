package debtors

import (
	"github.com/Pr3c10us/gbt/internals/app/debtors/command"
	"github.com/Pr3c10us/gbt/internals/app/debtors/query"
	"github.com/Pr3c10us/gbt/internals/domain/debtor"
)

type Commands struct {
	AddDebtor    command.AddDebtorCommand
	RemoveDebtor command.RemoveDebtorCommand
}

type Queries struct {
	GetDebtorByID query.GetDebtorByIDQuery
}

type Services struct {
	Commands Commands
	Queries  Queries
}

func NewDebtorServices(repository debtor.Repository) Services {
	return Services{
		Commands: Commands{
			AddDebtor:    command.NewAddDebtorCommand(repository),
			RemoveDebtor: command.NewRemoveDebtorCommand(repository),
		},
		Queries: Queries{
			GetDebtorByID: query.NewGetDebtorByID(repository),
		},
	}
}
