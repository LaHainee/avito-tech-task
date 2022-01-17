package models

type TransferRequest struct {
	SenderID   int64   `json:"sender_id,omitempty" form:"sender_id" validate:"required" example:"1"`
	ReceiverID int64   `json:"receiver_id,omitempty" form:"receiver_id" validate:"required" example:"2"`
	Amount     float64 `json:"amount,omitempty" form:"amount" validate:"required" example:"1000"`
}

type TransferUsersData struct {
	Sender   *UserData `json:"sender"`
	Receiver *UserData `json:"receiver"`
}
