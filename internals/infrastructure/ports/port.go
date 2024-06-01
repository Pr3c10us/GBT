package ports

import (
	"github.com/Pr3c10us/gbt/internals/app"
	"github.com/Pr3c10us/gbt/internals/infrastructure/ports/http"
	"github.com/Pr3c10us/gbt/pkg/logger"
)

type Port struct {
	Server http.Server
}

func NewPort(services app.Services, sugaredLogger logger.Logger) Port {
	return Port{
		Server: http.NewServer(services, sugaredLogger),
	}
}
