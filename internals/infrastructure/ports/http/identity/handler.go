package identity

import (
	identityService "github.com/Pr3c10us/gbt/internals/app/identity"
	identityEntity "github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/pkg/configs"
	"github.com/Pr3c10us/gbt/pkg/response"
	"github.com/Pr3c10us/gbt/pkg/utils"
	"github.com/Pr3c10us/gbt/pkg/validator"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service              identityService.Service
	environmentVariables *configs.EnvironmentVariables
}

func NewIdentityHandler(service identityService.Service, environmentVariables *configs.EnvironmentVariables) Handler {
	return Handler{
		service:              service,
		environmentVariables: environmentVariables,
	}
}

func (handler *Handler) Registration(context *gin.Context) {
	var addUserRequest identityEntity.AddUserRequest
	if err := context.ShouldBind(&addUserRequest); err != nil {
		err = validator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}

	err := handler.service.CreateUser.Handle(addUserRequest)
	if err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse(gin.H{"message": "user created"}, nil).Send(context)
}

// Authentication Struct
type authenticationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (handler *Handler) Authentication(context *gin.Context) {
	var authRequest authenticationRequest
	if err := context.ShouldBind(&authRequest); err != nil {
		err = validator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}

	user, err := handler.service.AuthenticateUser.Handle(authRequest.Username, authRequest.Password)
	if err != nil {
		_ = context.Error(err)
		return
	}

	token, err := utils.CreateTokenFromUser(user, handler.environmentVariables.JWTSecret)
	if err != nil {
		_ = context.Error(err)
		return
	}

	session := sessions.Default(context)
	session.Set("token", token)
	session.Options(sessions.Options{MaxAge: 60 * 60 * 24 * 31}) // 30 days
	if err := session.Save(); err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse(gin.H{"user": user, "token": token}, nil).Send(context)
}

func (handler *Handler) Logoff(context *gin.Context) {
	session := sessions.Default(context)
	session.Clear()
	if err := session.Save(); err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse(gin.H{"message": "logged off"}, nil).Send(context)
}
