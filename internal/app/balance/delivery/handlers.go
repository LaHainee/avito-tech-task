package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

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
	server.POST("/api/v1/transfer", h.Transfer)

	server.GET("/api/v1/balance/:user_id", h.GetBalance)
}

// Transfer
// @Summary 	Transfer money between users
// @Produce 	json
// @Param 		data body models.TransferRequest true "Data for transferring money"
// @Success 	200 {object} models.TransferUsersData
// @Failure		400 {object} models.ResponseMessage "Invalid request body"
// @Failure		404 {object} models.ResponseMessage "Sender not found | receiver not found"
// @Failure		422 {object} models.ResponseMessage "Not enough money"
// @Failure		500 {object} models.ResponseMessage "Internal server error"
// @Router 		/transfer [POST]
func (h *Handlers) Transfer(ctx echo.Context) error {
	h.logger.Info("Called handler Transfer for POST /api/v1/transfer")

	var transferData models.TransferRequest
	if err := ctx.Bind(&transferData); err != nil {
		h.logger.Warnf("Could not bind request body to models.RequestUpdateBalance: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			&models.ResponseMessage{Message: constants.InvalidBodyMessage})
	}
	h.logger.Infof("Request data: %v", transferData)

	transferResult, err := h.service.MakeTransfer(&transferData)
	if err != nil {
		switch errors.Is(err, createdErrors.ErrNotEnoughMoney) {
		case true:
			h.logger.Warnf("Unprocesseable request: %s", err)
			return ctx.JSON(
				http.StatusUnprocessableEntity,
				&models.ResponseMessage{Message: err.Error()})
		default:
			switch errors.Is(err, createdErrors.ErrSenderDoesNotExist) ||
				errors.Is(err, createdErrors.ErrReceiverDoesNotExist) {
			case true:
				h.logger.Warnf("%s", err)
				return ctx.JSON(
					http.StatusNotFound,
					&models.ResponseMessage{Message: err.Error()})
			default:
				h.logger.Errorf("Internal server error: %s", err)
				return ctx.JSON(
					http.StatusInternalServerError,
					&models.ResponseMessage{Message: err.Error()})
			}
		}
	}

	h.logger.Infof("Money transfer from %d to %d was successfully processed, received response: %v",
		transferData.SenderID, transferData.ReceiverID, transferResult)
	return ctx.JSON(http.StatusOK, transferResult)
}

// GetBalance
// @Summary 	Get user balance
// @Produce 	json
// @Param 		user_id path int true "User ID in BalanceApplication"
// @Param 		currency query string false "Currency to convert in"
// @Success 	200 {object} models.UserData
// @Failure		400 {object} models.ResponseMessage "Invalid user ID in query param"
// @Failure		404 {object} models.ResponseMessage "User not found"
// @Failure		422 {object} models.ResponseMessage "Unsupported currency"
// @Failure		500 {object} models.ResponseMessage "Internal server error"
// @Router 		/balance/{user_id} [GET]
func (h *Handlers) GetBalance(ctx echo.Context) error {
	h.logger.Info("Called handler GetBalance for GET /api/v1/balance/:id")

	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		h.logger.Warnf("Could not convert user id from string to int: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			&models.ResponseMessage{Message: constants.InvalidUserIDMessage})
	}
	currency := ctx.QueryParam("currency")
	h.logger.Infof("Request data: userID: %d, currency: %s", userID, currency)

	balance, err := h.service.GetBalance(userID, currency)
	switch errors.Is(err, createdErrors.ErrNotSupportedCurrency) {
	case true:
		h.logger.Warnf("Bad request: %s", err)
		return ctx.JSON(
			http.StatusUnprocessableEntity,
			&models.ResponseMessage{Message: err.Error()})
	case false:
		switch errors.Is(err, createdErrors.ErrUserDoesNotExist) {
		case true:
			h.logger.Warnf("%s", err)
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
	}

	h.logger.Infof("Request was successfully processed, received user data: %v", balance)
	return ctx.JSON(http.StatusOK, balance)
}

// UpdateBalance
// @Summary 	Update user balance
// @Produce 	json
// @Param 		user_id path int true "User ID in BalanceApplication"
// @Param 		data body models.RequestUpdateBalance true "Data for updating balance, operation = 0 - add money,operation = 1 - write off money"
// @Success 	200 {object} models.UserData
// @Failure		400 {object} models.ResponseMessage "Invalid user ID in query param | invalid request body"
// @Failure		422 {object} models.ResponseMessage "Not enough money | Not supported operation type | Amount field is required | Negative user ID"
// @Failure		500 {object} models.ResponseMessage "Internal server error"
// @Router 		/balance/{user_id} [POST]
func (h *Handlers) UpdateBalance(ctx echo.Context) error {
	h.logger.Info("Called handler UpdateBalance for POST /api/v1/balance/:id")

	var updateData models.RequestUpdateBalance
	if err := ctx.Bind(&updateData); err != nil {
		h.logger.Warnf("Could not bind request body to models.RequestUpdateBalance: %s", err)
		return ctx.JSON(
			http.StatusBadRequest,
			&models.ResponseMessage{Message: constants.InvalidBodyMessage})
	}
	h.logger.Infof("Request data: %v", updateData)

	userData, err := h.service.UpdateBalance(&updateData)
	if errors.Is(err, createdErrors.ErrNotEnoughMoney) || errors.Is(err, createdErrors.ErrNotSupportedOperationType) ||
		errors.Is(err, createdErrors.ErrAmountFiledIsRequired) || errors.Is(err, createdErrors.ErrNegativeUserID) {
		h.logger.Warnf("Bad request: %s", err)
		return ctx.JSON(
			http.StatusUnprocessableEntity,
			&models.ResponseMessage{Message: err.Error()})
	} else if err != nil {
		h.logger.Errorf("Internal server error: %s", err)
		return ctx.JSON(
			http.StatusInternalServerError,
			&models.ResponseMessage{Message: err.Error()})
	}

	h.logger.Infof("Request was successfully processed, received response: %v", userData)
	return ctx.JSON(http.StatusOK, userData)
}
