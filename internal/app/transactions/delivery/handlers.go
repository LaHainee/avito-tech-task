package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/app/transactions"
	"avito-tech-task/internal/pkg/constants"
	createdErrors "avito-tech-task/internal/pkg/errors"
)

type Handlers struct {
	service transactions.Service
	logger  *logrus.Logger
}

func NewHandlers(service transactions.Service, logger *logrus.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

func (h *Handlers) InitHandlers(server *echo.Echo) {
	server.GET("/api/v1/transactions/:user_id", h.GetTransactions)
}

// GetTransactions
// @Summary 	Get list of user transactions
// @Produce 	json
// @Param 		user_id path int true "User ID in BalanceApplication"
// @Param 		params body models.TransactionsSelectionParams true "Parameters for transactions selection"
// @Success 	200 {object} models.Transactions
// @Failure		400 {object} models.ResponseMessage "Invalid user ID in query param | invalid body"
// @Failure		404 {object} models.ResponseMessage "User not found"
// @Failure		500 {object} models.ResponseMessage "Internal server error"
// @Router 		/transactions/{user_id} [POST]
func (h *Handlers) GetTransactions(ctx echo.Context) error {
	h.logger.Info("Called handler GetTransactions for GET /api/v1/transactions/:user_id")

	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		h.logger.Warnf("Could not convert user id from string to int: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			&models.ResponseMessage{Message: constants.InvalidUserIDMessage})
	}

	var params models.TransactionsSelectionParams
	if err = ctx.Bind(&params); err != nil {
		h.logger.Warnf("Could not bind query params to models.TransactionsSelectionParams: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			&models.ResponseMessage{Message: constants.InvalidQueryParams})
	}

	transactions, err := h.service.GetUserTransactions(userID, &params)
	switch errors.Is(err, createdErrors.ErrUserDoesNotExist) {
	case true:
		h.logger.Warnf("Bad request: %s", err)
		return ctx.JSON(
			http.StatusNotFound,
			&models.ResponseMessage{Message: err.Error()})
	case false:
		if err != nil {
			h.logger.Errorf("Internal server error: %s", err)
			return ctx.JSON(
				http.StatusInternalServerError,
				&models.ResponseMessage{Message: err.Error()})
		}
	}

	h.logger.Infof("Request was successfully processed, received response: %v", transactions)
	return ctx.JSON(http.StatusOK, transactions)
}
