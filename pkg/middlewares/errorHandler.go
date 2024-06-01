package middlewares

import (
	"errors"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/Pr3c10us/gbt/pkg/logger"
	"github.com/Pr3c10us/gbt/pkg/response"
	"github.com/Pr3c10us/gbt/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
)

func ErrorHandlerMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			var (
				pqErr           *pq.Error
				customError     *appError.CustomError
				validationError *validator.ValidationError
			)
			logger.LogWithFields("error", "Error handler message", zap.Error(err))

			switch {
			case errors.As(err.Err, &pqErr):
				{
					switch pqErr.Code.Name() {
					case "unique_violation":
						response.ErrorResponse{
							StatusCode:   http.StatusConflict,
							Message:      "unique key value violated",
							ErrorMessage: pqErr.Detail,
						}.Send(c)
						return
					default:
						response.NewErrorResponse(pqErr).Send(c)
						return
					}
				}
			case errors.As(err.Err, &customError):
				response.ErrorResponse{
					StatusCode:   customError.StatusCode,
					Message:      customError.Message,
					ErrorMessage: customError.ErrorMessage,
				}.Send(c)
				return
			case errors.As(err.Err, &validationError):
				response.ErrorResponse{
					StatusCode:   validationError.StatusCode,
					Message:      validationError.Message,
					ErrorMessage: validationError.ErrorMessage,
				}.Send(c)
			default:
				response.NewErrorResponse(err).Send(c)
				return
			}
		}
	}
}
