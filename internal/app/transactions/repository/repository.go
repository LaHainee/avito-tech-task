package repository

import (
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/pkg/utils"
	"context"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	db utils.PgxIface
}

func NewStorage(pool utils.PgxIface) *Storage {
	return &Storage{pool}
}

func (s *Storage) GetUserTransactions(id int64, params *models.TransactionsSelectionParams) (models.Transactions, error) {
	var (
		rows pgx.Rows
		err  error
	)

	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	query := `SELECT description, amount, created FROM transactions WHERE user_id = $1 `

	if params.Since == "" { // no filter by transaction time
		if params.OrderAmount && params.OrderDate { // order by amount and time
			query += `ORDER BY amount DESC, created DESC LIMIT NULLIF($2, 0)`
		} else if params.OrderAmount { // order by amount
			query += `ORDER BY amount DESC LIMIT NULLIF($2, 0)`
		} else if params.OrderDate { // order by transaction time
			query += `ORDER BY created DESC LIMIT NULLIF($2, 0)`
		} else {
			query += `LIMIT NULLIF($2, 0)`
		}

		rows, err = transaction.Query(context.Background(), query, id, params.Limit)
		if err != nil {
			return nil, err
		}
	} else { // since transaction time
		if params.OrderAmount && params.OrderDate { // order by amount and time
			query += `AND created <= $2 ORDER BY amount DESC, created DESC LIMIT NULLIF($3, 0)`
		} else if params.OrderAmount { // order by amount
			query += `AND created <= $2 ORDER BY amount DESC LIMIT NULLIF($3, 0)`
		} else if params.OrderDate { // order by transaction time
			query += `AND created <= $2 ORDER BY created DESC LIMIT NULLIF($3, 0)`
		} else {
			query += `AND created <= $2 LIMIT NULLIF($3, 0)`
		}

		rows, err = transaction.Query(context.Background(), query, id, params.Since, params.Limit)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	userTransactions := models.Transactions{}
	for rows.Next() {
		var userTransaction models.Transaction
		if err = rows.Scan(&userTransaction.Description, &userTransaction.Amount, &userTransaction.Created); err != nil {
			return nil, err
		}
		userTransactions = append(userTransactions, &userTransaction)
	}

	return userTransactions, nil
}

func (s *Storage) DoesUserExist(id int64) (bool, error) {
	query := `SELECT user_id FROM balance WHERE user_id = $1`

	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	if err = transaction.QueryRow(context.Background(), query, id).Scan(&id); err != nil {
		if err != pgx.ErrNoRows {
			return false, err
		}
		err = nil
		return false, nil
	}
	return true, nil
}
