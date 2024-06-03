package http

import (
	"github.com/Pr3c10us/gbt/internals/app"
	"github.com/Pr3c10us/gbt/internals/infrastructure/ports/http/debtors"
	"github.com/Pr3c10us/gbt/internals/infrastructure/ports/http/identity"
	"github.com/Pr3c10us/gbt/pkg/configs"
	"github.com/Pr3c10us/gbt/pkg/logger"
	"github.com/Pr3c10us/gbt/pkg/middlewares"
	"github.com/Pr3c10us/gbt/pkg/response"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	environmentVariables = configs.LoadEnvironment()
	cookieStore          = cookie.NewStore([]byte(environmentVariables.CookieSecret))
)

type Server struct {
	services      app.Services
	engine        *gin.Engine
	sugaredLogger logger.Logger
}

func NewServer(services app.Services, sugaredLogger logger.Logger) Server {
	server := Server{
		engine:        gin.Default(),
		services:      services,
		sugaredLogger: sugaredLogger,
	}
	// Middlewares
	server.engine.Use(middlewares.RequestLoggingMiddleware(sugaredLogger))
	server.engine.Use(sessions.Sessions("gbt", cookieStore))
	server.engine.Use(middlewares.ErrorHandlerMiddleware(sugaredLogger))
	server.engine.NoRoute(middlewares.RouteNotFoundMiddleware())

	server.Health()
	server.Identity()
	server.Debtors()

	return server
}

func (server *Server) Health() {
	server.engine.GET("/health", func(c *gin.Context) {
		response.NewSuccessResponse(nil, nil).Send(c)
	})
}

func (server *Server) Identity() {
	handler := identity.NewIdentityHandler(server.services.IdentityService, environmentVariables)
	route := server.engine.Group("/api/v1/identity")
	{
		route.POST("/register", handler.Registration)
		route.POST("/authenticate", handler.Authentication)
		route.POST("/logoff", handler.Logoff)
	}
}

func (server *Server) Debtors() {
	handler := debtors.NewDebtorsHandler(server.services.DebtorServices)
	route := server.engine.Group("/api/v1/debtors")
	{
		route.POST("/", handler.AddDebtor)
		route.DELETE("/:id", handler.RemoveDebtor)
	}
}

func (server *Server) Run() {
	err := server.engine.Run(environmentVariables.Port)
	if err != nil {
		panic("Failed to start server")
	}
}
