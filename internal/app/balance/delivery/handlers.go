package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"avito-tech-task/internal/app/balance"
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/pkg/constants"
	createdErrors "avito-tech-task/internal/pkg/errors"
)

type Handlers struct {
	service balance.Service
	logger  *logrus.Logger
}

func NewHandlers(service balance.Service, logger *logrus.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger}
}

func (h *Handlers) InitHandlers(server *echo.Echo) {
	server.POST("/api/v1/balance/:user_id", h.UpdateBalance)
	server.GET("/api/v1/balance/:user_id", h.GetBalance)
}

func (h *Handlers) GetBalance(ctx echo.Context) error {
	h.logger.Info("Called handler GetBalance for GET /api/v1/balance/:id")


	userID, err := strconv.ParseInt(ctx.Param("user_id"), 32, 64)
	if err != nil {
		h.logger.Warnf("Could not convert user id from string to int: %s", err)
		return ctx.JSON(
			http.StatusUnprocessableEntity,
			&models.ResponseMessage{Message: constants.InvalidUserIDMessage})
	}

	userData, err := h.service.GetUserData(userID)
	if err != nil {
		switch err {
		case createdErrors.ErrUserDoesNotExist:
			h.logger.Warnf("User with id %d not found", userID)
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

	h.logger.Info("Request was successfully processed, received user data: %v", userData)
	return ctx.JSON(http.StatusOK, userData)
}

func (h *Handlers) UpdateBalance(ctx echo.Context) error {
	h.logger.Info("Called handler UpdateBalance for POST /api/v1/balance/:id")

	var request models.RequestUpdateBalance
	if err := ctx.Bind(&request); err != nil {
		h.logger.Warnf("Could not bind request body to models.RequestUpdateBalance: %s", err)
		return ctx.JSON(
			http.StatusUnprocessableEntity,
			&models.ResponseMessage{Message: constants.InvalidBodyMessage})
	}

	userData, err := h.service.UpdateBalance(&request)
	if err != nil {
		switch err {
		case createdErrors.ErrNotEnoughMoney, createdErrors.ErrNotSupportedOperationType,
			createdErrors.ErrAmountFiledIsRequired, createdErrors.ErrNegativeUserID:
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

	h.logger.Info("Request was successfully processed, received updated user data: %v", userData)
	return ctx.JSON(http.StatusOK, userData)
}
