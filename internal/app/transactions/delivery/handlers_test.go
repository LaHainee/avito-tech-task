package delivery

import (
	"avito-tech-task/config"
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/app/transactions/mock"
	"avito-tech-task/internal/pkg/constants"
	createdErrors "avito-tech-task/internal/pkg/errors"
	"avito-tech-task/internal/pkg/utils"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHandlers_GetTransactions(t *testing.T) {
	const removeLogs = true // set false to deny deleting logs after test

	config := &config.Config{
		LoggingLevel:    "debug",
		LoggingFilePath: "./logs/",
		CurrencyApiURL:  "",
		Server:          config.ServerConfig{},
	}
	logger, fileClose := utils.NewLogger(config)
	defer func(close func() error) {
		if err := close(); err != nil {
			t.Errorf("Could not close file: %s", err)
		}
	}(fileClose)
	internalServerErr := errors.New("Internal server error")

	if removeLogs {
		defer func() {
			if err := os.RemoveAll("./logs/"); err != nil {
				t.Errorf("Could not remove temporary logs directory: %s", err)
			}
		}()
	}

	timeNow := time.Now()
	tests := []struct {
		name           string
		serviceMock    *mock.MockService
		userIDParam    string
		body           string
		expectedStatus int
		expected       interface{}
	}{
		{
			name: "Successfully get user transactions list",
			serviceMock: &mock.MockService{
				GetUserTransactionsFunc: func(n int64, transactionsSelectionParams *models.TransactionsSelectionParams) (models.Transactions, error) {
					return models.Transactions{
						&models.Transaction{
							Description: "description",
							Amount:      1000,
							Created:     timeNow,
						},
					}, nil
				},
			},
			userIDParam:    "1",
			body:           `{"limit": 10}`,
			expectedStatus: http.StatusOK,
			expected: models.Transactions{
				&models.Transaction{
					Description: "description",
					Amount:      1000,
					Created:     timeNow,
				},
			},
		},
		{
			name: "User does not exist",
			serviceMock: &mock.MockService{
				GetUserTransactionsFunc: func(n int64, transactionsSelectionParams *models.TransactionsSelectionParams) (models.Transactions, error) {
					return nil, createdErrors.ErrUserDoesNotExist
				},
			},
			userIDParam:    "1",
			body:           `{"limit": 10}`,
			expectedStatus: http.StatusNotFound,
			expected:       &models.ResponseMessage{Message: createdErrors.ErrUserDoesNotExist.Error()},
		},
		{
			name: "Internal server error",
			serviceMock: &mock.MockService{
				GetUserTransactionsFunc: func(n int64, transactionsSelectionParams *models.TransactionsSelectionParams) (models.Transactions, error) {
					return nil, internalServerErr
				},
			},
			userIDParam:    "1",
			body:           `{"limit": 10}`,
			expectedStatus: http.StatusInternalServerError,
			expected:       &models.ResponseMessage{Message: internalServerErr.Error()},
		},
		{
			name:           "Invalid user ID as param",
			userIDParam:    "hello",
			expectedStatus: http.StatusBadRequest,
			expected:       &models.ResponseMessage{Message: constants.InvalidUserIDMessage},
		},
		{
			name:           "Invalid query params",
			userIDParam:    "1",
			body:           `{"limit": "string???"}`,
			expectedStatus: http.StatusBadRequest,
			expected:       &models.ResponseMessage{Message: constants.InvalidQueryParams},
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			server := echo.New()

			req := httptest.NewRequest(echo.POST, "/", strings.NewReader(test.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/transactions/:user_id")
			ctx.SetParamNames("user_id")
			ctx.SetParamValues(test.userIDParam)

			handlers := NewHandlers(test.serviceMock, logger)
			if assert.NoError(t, handlers.GetTransactions(ctx)) {
				assert.Equal(t, test.expectedStatus, rec.Code)

				expectedString, _ := json.Marshal(test.expected)
				assert.Equal(t, string(expectedString)+"\n", rec.Body.String())
			}
		})
	}
}
