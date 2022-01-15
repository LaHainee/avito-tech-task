package models

type UserData struct {
	UserID  int64   `json:"user_id,omitempty"`
	Balance float64 `json:"balance,omitempty"`
}
