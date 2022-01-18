package models

import "time"

type Transaction struct {
	OperationType string    `json:"operation_type"`
	ReceiverID    int64     `json:"receiver_id,omitempty"`
	Amount        float64   `json:"amount"`
	Created       time.Time `json:"created"`
}

type TransactionsSelectionParams struct {
	Limit         int    `json:"limit,omitempty" form:"limit"`
	Since         string `json:"since,omitempty" form:"since"`
	OperationType int    `json:"operation_type,omitempty" form:"operation_type"`
	OrderAmount   bool   `json:"order_amount,omitempty" form:"order_amount"`
	OrderDate     bool   `json:"order_date,omitempty" form:"order_date"`
}

type Transactions []*Transaction
