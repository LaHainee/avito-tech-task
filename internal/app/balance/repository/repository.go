package repository

import (
	"avito-tech-task/internal/app/models"
	"fmt"
	"github.com/jackc/pgx"
)

const (
	queryUpdateBalance     = `UPDATE balance SET balance = balance + $1 WHERE user_id = $2 RETURNING balance`
	queryAddNewTransaction = `INSERT INTO transactions(description, amount, user_id) VALUES ($1, $2, $3)`
	queryGetUserBalance    = `SELECT balance FROM balance WHERE user_id = $1`
	queryAddNewAccount     = `INSERT INTO balance (user_id, balance) VALUES($1, 0)`
)

type Storage struct {
	pool *pgx.ConnPool
}

func NewStorage(pool *pgx.ConnPool) *Storage {
	return &Storage{pool}
}

func (s *Storage) CreateAccount(id int64) error {
	_, err := s.pool.Exec(queryAddNewAccount, id)

	return err
}

func (s *Storage) GetUserData(id int64) (*models.UserData, error) {
	var balance float64

	if err := s.pool.QueryRow(queryGetUserBalance, id).Scan(&balance); err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}

	return &models.UserData{UserID: id, Balance: balance}, nil
}

func (s *Storage) MakeTransfer(senderID, receiverID int64, amount float64) error {
	transaction, err := s.pool.Begin() // start transactions for safe money transfer
	defer func() {
		if err != nil {
			_ = transaction.Rollback()
		} else {
			_ = transaction.Commit()
		}
	}()

	if _, err = transaction.Exec(queryUpdateBalance, amount*-1, senderID); err != nil {
		return err
	}
	if _, err = transaction.Exec(queryUpdateBalance, amount, receiverID); err != nil {
		return err
	}
	if _, err = transaction.Exec(queryAddNewTransaction, fmt.Sprintf("Sent %.2fRUB to user %d", amount,
		receiverID), amount, senderID); err != nil {
		return err
	}
	if _, err = transaction.Exec(queryAddNewTransaction, fmt.Sprintf("Recevied %.2fRUB from user %d", amount,
		senderID), amount, receiverID); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateBalance(id int64, amount float64) (float64, error) {
	transaction, err := s.pool.Begin()
	defer func() {
		if err != nil {
			_ = transaction.Rollback()
		} else {
			_ = transaction.Commit()
		}
	}()

	var balance float64
	if err = transaction.QueryRow(queryUpdateBalance, amount, id).Scan(&balance); err != nil {
		return 0, err
	}

	var operationDescription string
	if amount < 0 {
		operationDescription = fmt.Sprintf("Write off %.2fRUB", amount*-1)
		amount *= -1
	} else {
		operationDescription = fmt.Sprintf("Add %.2fRUB", amount)
	}

	if _, err = transaction.Exec(queryAddNewTransaction, operationDescription, amount, id); err != nil {
		return 0, err
	}

	return balance, nil
}
