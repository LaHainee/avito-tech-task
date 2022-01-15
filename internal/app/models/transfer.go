package models

type TransferRequest struct {
	SenderID   int64   `json:"sender_id,omitempty" form:"sender_id" validate:"required"`
	ReceiverID int64   `json:"receiver_id,omitempty" form:"receiver_id" validate:"required"`
	Amount     float64 `json:"amount,omitempty" form:"amount" validate:"required"`
}

type TransferResponse struct {
	Sender   UserData `json:"sender"`
	Receiver UserData `json:"receiver"`
}
