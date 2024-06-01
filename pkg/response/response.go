package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	StatusCode   int    `json:"statusCode"`
	Message      string `json:"message"`
	ErrorMessage any    `json:"error"`
}

func NewErrorResponse(err error) ErrorResponse {

	var errorMessage string = ""
	if err != nil {
		errorMessage = err.Error()
	}

	return ErrorResponse{
		StatusCode:   http.StatusInternalServerError,
		Message:      "Internal Server Error",
		ErrorMessage: errorMessage,
	}
}

func (res ErrorResponse) Send(c *gin.Context) {
	c.AbortWithStatusJSON(res.StatusCode, res)
}

type SuccessResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`     // data payload
	Metadata   interface{} `json:"metadata,omitempty"` //pagination and other payload
}

func NewSuccessResponse(data interface{}, metadata interface{}) *SuccessResponse {
	return &SuccessResponse{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Data:       data,
		Metadata:   metadata,
	}
}

func (res SuccessResponse) Send(c *gin.Context) {
	c.JSON(res.StatusCode, res)
}
