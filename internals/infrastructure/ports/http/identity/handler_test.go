package identity

import (
	"encoding/json"
	"errors"
	identityService "github.com/Pr3c10us/gbt/internals/app/identity"
	"github.com/Pr3c10us/gbt/internals/app/identity/queries"
	identityRepository "github.com/Pr3c10us/gbt/internals/domain/identity"
	"github.com/Pr3c10us/gbt/pkg/appError"
	"github.com/Pr3c10us/gbt/pkg/configs"
	"github.com/Pr3c10us/gbt/pkg/logger"
	"github.com/Pr3c10us/gbt/pkg/middlewares"
	"github.com/Pr3c10us/gbt/pkg/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	environmentVariables = configs.LoadEnvironment()
	cookieStore          = cookie.NewStore([]byte(environmentVariables.CookieSecret))
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ErrorHandlerMiddleware(logger.NewSugarLogger(false)))
	r.Use(sessions.Sessions("gbt", cookieStore))
	return r
}

func TestHandler_Registration(t *testing.T) {
	tests := []struct {
		name        string
		args        identityRepository.AddUserRequest
		statusCode  int
		ExpectedErr error
	}{
		{
			name: "valid response",
			args: identityRepository.AddUserRequest{
				Username:         "MrMan",
				Password:         "secretMan",
				SecurityQuestion: "who strong",
				SecurityAnswer:   "me",
			},
			statusCode:  http.StatusOK,
			ExpectedErr: nil,
		},
		{
			name: "invalid body",
			args: identityRepository.AddUserRequest{
				Username:         "",
				Password:         "",
				SecurityQuestion: "",
				SecurityAnswer:   "",
			},
			statusCode:  http.StatusNotAcceptable,
			ExpectedErr: nil,
		},
		{
			name: "duplicate username",
			args: identityRepository.AddUserRequest{
				Username:         "MrMan",
				Password:         "secretMan",
				SecurityQuestion: "who strong",
				SecurityAnswer:   "me",
			},
			statusCode: http.StatusConflict,
			ExpectedErr: &pq.Error{
				Code:     "23505",
				Severity: "ERROR",
				Message:  "duplicate key value violates unique constraint \"username\"",
				Detail:   "Key (username)=(MrMan) already exists.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser, _ := identityRepository.NewUser(tt.args)

			mockRepo := new(identityRepository.MockRepository)
			mockRepo.On("CreateUser", mock.MatchedBy(utils.UserMatcher(mockUser))).Return(tt.ExpectedErr)

			service := identityService.NewIdentityService(mockRepo)
			handler := NewIdentityHandler(service, environmentVariables)

			engine := setupRouter()
			engine.POST("/api/v1/identity/register", handler.Registration)

			w := httptest.NewRecorder()
			addUserRequestJSON, _ := json.Marshal(tt.args)
			r, _ := http.NewRequest("POST", "/api/v1/identity/register", strings.NewReader(string(addUserRequestJSON)))
			r.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(w, r)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}

}

func TestHandler_Authentication(t *testing.T) {
	pw := "secretMan"
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		t.Error(err)
	}

	user := identityRepository.User{
		ID:               uuid.New().String(),
		Username:         "MrMan",
		Password:         string(hashedPw),
		SecurityQuestion: "who strong",
		SecurityAnswer:   "me",
	}
	token, err := utils.CreateTokenFromUser(&user, "secret")
	if err != nil {
		t.Error(err)
	}

	type want struct {
		user  identityRepository.User
		token string
	}
	tests := []struct {
		name        string
		args        authenticationRequest
		want        want
		statusCode  int
		ExpectedErr error
	}{
		{
			name: "valid credentials",
			args: authenticationRequest{
				Username: "MrMan",
				Password: pw,
			},
			want: want{
				user:  user,
				token: token,
			},
			statusCode:  http.StatusOK,
			ExpectedErr: nil,
		},
		{
			name: "invalid credentials",
			args: authenticationRequest{
				Username: "MrMan2",
				Password: pw,
			},
			want: want{
				user:  identityRepository.User{},
				token: token,
			},
			statusCode:  http.StatusNotFound,
			ExpectedErr: appError.NotFound(errors.New("user with username MrMan2 does not exit")),
		},
		{
			name: "incorrect password",
			args: authenticationRequest{
				Username: "MrMan",
				Password: "wrongPass",
			},
			want: want{
				user:  user,
				token: token,
			},
			statusCode:  http.StatusBadRequest,
			ExpectedErr: appError.BadRequest(queries.IncorrectPasswordErr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(identityRepository.MockRepository)
			mockRepo.On("GetUser", tt.args.Username).Return(tt.want.user, tt.ExpectedErr)

			service := identityService.NewIdentityService(mockRepo)
			handler := NewIdentityHandler(service, environmentVariables)

			engine := setupRouter()
			engine.POST("/api/v1/identity/authenticate", handler.Authentication)

			w := httptest.NewRecorder()
			authReq, _ := json.Marshal(tt.args)
			r, _ := http.NewRequest("POST", "/api/v1/identity/authenticate", strings.NewReader(string(authReq)))
			r.Header.Set("Content-Type", "application/json")

			engine.ServeHTTP(w, r)
			log.Print(w.Body.String())
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
