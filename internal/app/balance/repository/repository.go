package repository

import (
	"avito-tech-task/internal/app/models"
	"github.com/jackc/pgx"
)

type Storage struct {
	pool *pgx.ConnPool
}

func NewStorage(pool *pgx.ConnPool) *Storage {
	return &Storage{pool}
}

func (s *Storage) CreateAccount(id int64) error {
	query := `INSERT INTO balance (user_id, balance) VALUES($1, 0)`

	_, err := s.pool.Exec(query, id)

	return err
}

func (s *Storage) GetUserData(id int64) (*models.UserData, error) {
	query := `SELECT balance FROM balance WHERE user_id = $1`

	var balance float64

	if err := s.pool.QueryRow(query, id).Scan(&balance); err != nil {
		if err != pgx.ErrNoRows {
			return nil, err
		}
		return nil, nil
	}

	return &models.UserData{UserID: id, Balance: balance}, nil
}

func (s *Storage) UpdateBalance(id int64, amount float64) (*models.UserData, error) {
	transaction, err := s.pool.Begin()
	defer func() {
		if err != nil {
			_ = transaction.Rollback()
		} else {
			_ = transaction.Commit()
		}
	}()

	query := `UPDATE balance SET balance = balance + $1 WHERE user_id = $2 RETURNING balance`

	var userData models.UserData
	if err = transaction.QueryRow(query, amount, id).Scan(&userData.Balance); err != nil {
		return nil, err
	}
	userData.UserID = id

	return &userData, nil
}
