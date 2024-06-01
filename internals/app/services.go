package app

import (
	"github.com/Pr3c10us/gbt/internals/app/identity"
	"github.com/Pr3c10us/gbt/internals/infrastructure/adapters"
)

type Services struct {
	IdentityService identity.Service
}

func NewServices(adapter adapters.Adapter) Services {
	return Services{
		IdentityService: identity.NewIdentityService(adapter.IdentityRepository),
	}
}
