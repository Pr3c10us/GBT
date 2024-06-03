package app

import (
	"github.com/Pr3c10us/gbt/internals/app/debtors"
	"github.com/Pr3c10us/gbt/internals/app/identity"
	"github.com/Pr3c10us/gbt/internals/infrastructure/adapters"
)

type Services struct {
	IdentityService identity.Service
	DebtorServices  debtors.Services
}

func NewServices(adapter adapters.Adapter) Services {
	return Services{
		IdentityService: identity.NewIdentityService(adapter.IdentityRepository),
		DebtorServices:  debtors.NewDebtorServices(adapter.DebtorsRepository),
	}
}
