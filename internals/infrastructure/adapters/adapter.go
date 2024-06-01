package adapters

import (
	"database/sql"
	"github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/internals/infrastructure/adapters/persistence/postgresql"
	"github.com/Pr3c10us/gbt/pkg/logger"
)

type Adapter struct {
	IdentityRepository identity.Repository
}

func NewAdapter(db *sql.DB, logger logger.Logger) Adapter {
	return Adapter{
		IdentityRepository: postgresql.NewIdentityRepositoryPG(db, logger),
	}
}
