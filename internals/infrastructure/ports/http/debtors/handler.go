package debtors

import (
	"github.com/Pr3c10us/gbt/internals/app/debtors"
	"github.com/Pr3c10us/gbt/internals/app/debtors/command"
	"github.com/Pr3c10us/gbt/pkg/response"
	appValidator "github.com/Pr3c10us/gbt/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DebtorID struct {
	ID string `json:"id" validate:"required,uuid"`
}

type Handler struct {
	services debtors.Services
}

func NewDebtorsHandler(services debtors.Services) Handler {
	return Handler{
		services: services,
	}
}

func (handler *Handler) AddDebtor(context *gin.Context) {
	var addDebtorReq command.AddDebtorRequest
	if err := context.ShouldBind(&addDebtorReq); err != nil {
		err = appValidator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}

	err := handler.services.Commands.AddDebtor.Handle(addDebtorReq)
	if err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse(gin.H{"message": "debtor added successfully"}, nil).Send(context)
}

func (handler *Handler) RemoveDebtor(context *gin.Context) {
	debtorID := DebtorID{
		ID: context.Param("id"),
	}
	validate := validator.New()
	if err := validate.Struct(debtorID); err != nil {
		err = appValidator.ValidateRequest(err)
		_ = context.Error(err)
		return
	}

	err := handler.services.Commands.RemoveDebtor.Handle(debtorID.ID)
	if err != nil {
		_ = context.Error(err)
		return
	}

	response.NewSuccessResponse(gin.H{"message": "debtor removed successfully"}, nil).Send(context)
}
