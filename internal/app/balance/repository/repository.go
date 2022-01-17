package repository

import (
	createdErrors "avito-tech-task/internal/pkg/errors"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/pkg/utils"
)

type Storage struct {
	db utils.PgxIface
}

func NewStorage(conn utils.PgxIface) *Storage {
	return &Storage{conn}
}

const (
	queryUpdateBalance   = `UPDATE balance SET balance = balance + $1 WHERE user_id = $2 RETURNING balance`
	querySaveTransaction = `INSERT INTO transactions(description, amount, user_id) VALUES ($1, $2, $3)`
	queryGetBalance      = `SELECT balance FROM balance WHERE user_id = $1`
	queryInsertBalance   = `INSERT INTO balance (user_id, balance) VALUES($1, 0)`
	queryGetUser         = `SELECT user_id, balance FROM balance WHERE user_id = $1`
)

func (s *Storage) CreateAccount(userID int64) error {
	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	_, err = transaction.Exec(context.Background(), queryInsertBalance, userID)

	return err
}

func (s *Storage) GetUserData(userID int64) (*models.UserData, error) {
	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	var balance float64
	if err = transaction.QueryRow(context.Background(), queryGetBalance, userID).Scan(&balance); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, createdErrors.ErrUserDoesNotExist
	}

	return &models.UserData{UserID: userID, Balance: balance}, nil
}

func (s *Storage) GetTransferUsersData(senderID, receiverID int64) (*models.TransferUsersData, error) {
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
	if err = transaction.QueryRow(context.Background(), queryGetUser, senderID).Scan(&transferUsers.Sender.UserID,
		&transferUsers.Sender.Balance); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		transferUsers.Sender = nil
		err = nil
	}

	// getting info about receiver
	if err = transaction.QueryRow(context.Background(), queryGetUser, receiverID).Scan(&transferUsers.Receiver.UserID,
		&transferUsers.Receiver.Balance); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		transferUsers.Receiver = nil
		err = nil
	}

	return transferUsers, nil
}

func (s *Storage) MakeTransfer(senderID, receiverID int64, amount float64) error {
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

func (s *Storage) UpdateBalance(userID int64, amount float64) (float64, error) {
	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	var balance float64
	if err = transaction.QueryRow(context.Background(), queryUpdateBalance, amount, userID).Scan(&balance); err != nil {
		return 0, err
	}

	var operationDescription string
	if amount < 0 {
		operationDescription = fmt.Sprintf("Write off %.2fRUB", amount*-1)
		amount *= -1
	} else {
		operationDescription = fmt.Sprintf("Add %.2fRUB", amount)
	}

	if _, err = transaction.Exec(context.Background(), querySaveTransaction, operationDescription, amount, userID); err != nil {
		return 0, err
	}

	return balance, nil
}
