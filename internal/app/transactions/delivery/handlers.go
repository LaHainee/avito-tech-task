package delivery

import (
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/app/transactions"
	"avito-tech-task/internal/pkg/constants"
	createdErrors "avito-tech-task/internal/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
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

func (h *Handlers) GetTransactions(ctx echo.Context) error {
	h.logger.Info("Called handler GetTransactions for GET /api/v1/transactions/:user_id")

	userID, err := strconv.ParseInt(ctx.Param("user_id"), 32, 64)
	if err != nil {
		h.logger.Warnf("Could not convert user id from string to int: %s", err)
		return ctx.JSON(
			http.StatusUnprocessableEntity,
			&models.ResponseMessage{Message: constants.InvalidUserIDMessage})
	}

	var params models.TransactionsSelectionParams
	if err = ctx.Bind(&params); err != nil {
		h.logger.Warnf("Could not bind query params to models.TransactionsSelectionParams: %s", err)
		return ctx.JSON(
			http.StatusUnprocessableEntity,
			&models.ResponseMessage{Message: constants.InvalidQueryParams})
	}

	log.Println(params)

	transactions, err := h.service.GetUserTransactions(userID, &params)
	if err != nil {
		switch err {
		case createdErrors.ErrUserDoesNotExist:
			h.logger.Warnf("Bad request: %s", err)
			return ctx.JSON(
				http.StatusBadRequest,
				&models.ResponseMessage{Message: err.Error()})
		default:
			h.logger.Errorf("Internal server error: %s", err)
			return ctx.JSON(
				http.StatusInternalServerError,
				&models.ResponseMessage{Message: err.Error()})
		}
	}

	h.logger.Info("Request was successfully processed, received response: %v", transactions)
	return ctx.JSON(http.StatusOK, transactions)
}
