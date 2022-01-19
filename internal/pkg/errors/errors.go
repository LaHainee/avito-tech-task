package errors

import "errors"

var (
	ErrNegativeUserID            = errors.New("user id must be positive integer")
	ErrNotEnoughMoney            = errors.New("not enough money on balance")
	ErrNotSupportedOperationType = errors.New("not supported operation type")
	ErrAmountFiledIsRequired     = errors.New("amount field is required and must be greater than zero")
	ErrUserDoesNotExist          = errors.New("user does not exist")
	ErrSenderDoesNotExist        = errors.New("sender does not exist")
	ErrReceiverDoesNotExist      = errors.New("receiver does not exist")
	ErrSenderIDisRequired        = errors.New("sender_id is required")
	ErrReceiverIDisRequired      = errors.New("receiver_id is required")
	ErrNotSupportedCurrency      = errors.New("currency is not supported")
	ErrNegativeLimit             = errors.New("limit value must be positive integer")
)
