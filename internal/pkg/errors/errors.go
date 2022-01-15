package errors

import "errors"

var (
	ErrNegativeUserID            = errors.New("user id must be positive integer")
	ErrNotEnoughMoney            = errors.New("not enough money on balance")
	ErrNotSupportedOperationType = errors.New("not supported operation type")
	ErrAmountFiledIsRequired     = errors.New("amount field is required")
	ErrUserDoesNotExist          = errors.New("user does not exist")
)
