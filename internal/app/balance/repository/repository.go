package repository

import (
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/pkg/utils"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	db utils.PgxIface
}

func NewStorage(conn utils.PgxIface) *Storage {
	return &Storage{conn}
}

func (s *Storage) CreateAccount(id int64) error {
	query := `INSERT INTO balance (user_id, balance) VALUES($1, 0)`

	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	_, err = transaction.Exec(context.Background(), query, id)

	return err
}

func (s *Storage) GetUserData(id int64) (*models.UserData, error) {
	query := `SELECT balance FROM balance WHERE user_id = $1`

	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	var balance float64
	if err = transaction.QueryRow(context.Background(), query, id).Scan(&balance); err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}

	return &models.UserData{UserID: id, Balance: balance}, nil
}

func (s *Storage) GetTransferUsersData(senderID, receiverID int64) (*models.TransferUsersData, error) {
	query := `SELECT user_id, balance FROM balance WHERE user_id = $1`

	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	transferUsers := &models.TransferUsersData{}
	transferUsers.Sender = &models.UserData{}
	transferUsers.Receiver = &models.UserData{}

	// getting info about sender
	if err = transaction.QueryRow(context.Background(), query, senderID).Scan(&transferUsers.Sender.UserID,
		&transferUsers.Sender.Balance); err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		transferUsers.Sender = nil
		err = nil
	}

	// getting info about receiver
	if err = transaction.QueryRow(context.Background(), query, receiverID).Scan(&transferUsers.Receiver.UserID,
		&transferUsers.Receiver.Balance); err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		transferUsers.Receiver = nil
		err = nil
	}

	return transferUsers, nil
}

func (s *Storage) MakeTransfer(senderID, receiverID int64, amount float64) error {
	queryUpdateBalance := `UPDATE balance SET balance = balance + $1 WHERE user_id = $2 RETURNING balance`
	querySaveTransaction := `INSERT INTO transactions(description, amount, user_id) VALUES ($1, $2, $3)`

	transaction, err := s.db.Begin(context.Background()) // start transactions for safe money transfer
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	if _, err = transaction.Exec(context.Background(), queryUpdateBalance, amount*-1, senderID); err != nil {
		return err
	}
	if _, err = transaction.Exec(context.Background(), queryUpdateBalance, amount, receiverID); err != nil {
		return err
	}
	if _, err = transaction.Exec(context.Background(), querySaveTransaction, fmt.Sprintf("Sent %.2fRUB to user %d", amount,
		receiverID), amount, senderID); err != nil {
		return err
	}
	if _, err = transaction.Exec(context.Background(), querySaveTransaction, fmt.Sprintf("Recevied %.2fRUB from user %d", amount,
		senderID), amount, receiverID); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateBalance(id int64, amount float64) (float64, error) {
	queryUpdateBalance := `UPDATE balance SET balance = balance + $1 WHERE user_id = $2 RETURNING balance`
	querySaveTransaction := `INSERT INTO transactions(description, amount, user_id) VALUES ($1, $2, $3)`

	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	var balance float64
	if err = transaction.QueryRow(context.Background(), queryUpdateBalance, amount, id).Scan(&balance); err != nil {
		return 0, err
	}

	var operationDescription string
	if amount < 0 {
		operationDescription = fmt.Sprintf("Write off %.2fRUB", amount*-1)
		amount *= -1
	} else {
		operationDescription = fmt.Sprintf("Add %.2fRUB", amount)
	}

	if _, err = transaction.Exec(context.Background(), querySaveTransaction, operationDescription, amount, id); err != nil {
		return 0, err
	}

	return balance, nil
}
