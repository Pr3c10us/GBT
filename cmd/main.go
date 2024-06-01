package main

import (
	"database/sql"
	"fmt"
	"github.com/Pr3c10us/gbt/internals/app"
	"github.com/Pr3c10us/gbt/internals/infrastructure/adapters"
	"github.com/Pr3c10us/gbt/internals/infrastructure/ports"
	"github.com/Pr3c10us/gbt/pkg/configs"
	"github.com/Pr3c10us/gbt/pkg/logger"
)

var (
	environmentVariables = configs.LoadEnvironment()
)

func main() {

	// PG_DB instantiation
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		environmentVariables.PGDB.Host,
		environmentVariables.PGDB.Port,
		environmentVariables.PGDB.Username,
		environmentVariables.PGDB.Password,
		environmentVariables.PGDB.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic("Failed to instantiation DB connection")
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	err = db.Ping()
	if err != nil {
		panic("No connection could be made because the target machine actively refused it")
	}

	sugaredLogger := logger.NewSugarLogger(environmentVariables.ProductionEnvironment)

	adapter := adapters.NewAdapter(db, sugaredLogger)
	services := app.NewServices(adapter)
	port := ports.NewPort(services, sugaredLogger)
	port.Server.Run()
}
