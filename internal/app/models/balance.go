package models

const (
	ADD = iota
	REDUCE
)

type RequestUpdateBalance struct {
	UserID        int64   `json:"user_id,omitempty" param:"user_id" validate:"gt=0"`
	OperationType int     `json:"operation_type,omitempty" form:"operation_type" validate:"operation_type"`
	Amount        float64 `json:"amount,omitempty" form:"amount" validate:"required"`
}

type UserData struct {
	UserID  int64   `json:"user_id,omitempty"`
	Balance float64 `json:"balance,omitempty"`
}
