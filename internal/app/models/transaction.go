package models

import "time"

type Transaction struct {
	Description string    `json:"description,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	Created     time.Time `json:"created"`
}

type TransactionsSelectionParams struct {
	Limit       int    `json:"limit,omitempty" query:"limit"`
	Since       string `json:"since,omitempty" query:"since"`
	OrderAmount bool   `json:"order_amount,omitempty" query:"order_amount"`
	OrderDate   bool   `json:"order_date,omitempty" query:"order_date"`
}

type Transactions []*Transaction
