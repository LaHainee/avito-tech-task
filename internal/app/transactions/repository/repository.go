package repository

import (
	"avito-tech-task/internal/app/models"
	"github.com/jackc/pgx"
)

const (
	queryCheckIfUserExist = `SELECT user_id FROM balance WHERE user_id = $1`
)

type Storage struct {
	pool *pgx.ConnPool
}

func NewStorage(pool *pgx.ConnPool) *Storage {
	return &Storage{pool}
}

func (s *Storage) GetUserTransactions(id int64, params *models.TransactionsSelectionParams) (models.Transactions, error) {
	var (
		rows *pgx.Rows
		err  error
	)
	query := `SELECT description, amount, created FROM transactions WHERE user_id = $1 `

	if params.Since == "" { // no filter by transaction time
		if params.OrderAmount && params.OrderDate {
			query += `ORDER BY amount DESC, created DESC LIMIT NULLIF($2, 0)`
		}

		rows, err = s.pool.Query(query, id, params.Limit)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	transactions := models.Transactions{}
	for rows.Next() {
		var transaction models.Transaction
		if err = rows.Scan(&transaction.Description, &transaction.Amount, &transaction.Created); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (s *Storage) DoesUserExist(id int64) (bool, error) {
	if err := s.pool.QueryRow(queryCheckIfUserExist, id).Scan(&id); err != nil {
		if err != pgx.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return true, nil
}
