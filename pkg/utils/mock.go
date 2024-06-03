package utils

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Pr3c10us/gbt/pkg/configs"
	"github.com/Pr3c10us/gbt/pkg/logger"
	"github.com/Pr3c10us/gbt/pkg/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	sugarLogger          = logger.NewSugarLogger(false)
	environmentVariables = configs.LoadEnvironment()
	cookieStore          = cookie.NewStore([]byte(environmentVariables.CookieSecret))
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		message := fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err)
		sugarLogger.Log("fatal", message)
	}

	return db, mock
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ErrorHandlerMiddleware(logger.NewSugarLogger(false)))
	r.Use(sessions.Sessions("gbt", cookieStore))
	return r
}
