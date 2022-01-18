package repository

import (
	"avito-tech-task/internal/pkg/constants"
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4"

	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/pkg/utils"
)

type Storage struct {
	db utils.PgxIface
}

func NewStorage(pool utils.PgxIface) *Storage {
	return &Storage{pool}
}

const (
	queryGetUserID = `SELECT user_id FROM balance WHERE user_id = $1`
)

//nolint:cyclop
func (s *Storage) GetUserTransactions(userID int64, params *models.TransactionsSelectionParams) (models.Transactions, error) {
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

	query := `SELECT operation_type, receiver, amount, created FROM transactions WHERE sender = $1 `

	switch params.OperationType {
	case constants.ADD:
		query += `AND operation_type = 'add' `
	case constants.REDUCE:
		query += `AND operation_type = 'write_off' `
	case constants.TRANSFER:
		query += `AND operation_type = 'transfer' `
	}

	if params.Since == "" { // no filter by transaction time
		switch params.OrderDate {
		case true:
			switch params.OrderAmount {
			case true:
				query += `ORDER BY amount DESC, created DESC LIMIT NULLIF($2, 0)`
			case false:
				query += `ORDER BY created DESC LIMIT NULLIF($2, 0)`
			}
		case false:
			switch params.OrderAmount {
			case true:
				query += `ORDER BY amount DESC LIMIT NULLIF($2, 0)`
			case false:
				query += `LIMIT NULLIF($2, 0)`
			}
		}
		rows, err = transaction.Query(context.Background(), query, userID, params.Limit)
		if err != nil {
			return nil, err
		}
	} else { // since transaction time
		switch params.OrderDate {
		case true:
			switch params.OrderAmount {
			case true:
				query += `AND created <= $2 ORDER BY amount DESC, created DESC LIMIT NULLIF($3, 0)`
			case false:
				query += `AND created <= $2 ORDER BY created DESC LIMIT NULLIF($3, 0)`
			}
		case false:
			switch params.OrderAmount {
			case true:
				query += `AND created <= $2 ORDER BY amount DESC LIMIT NULLIF($3, 0)`
			case false:
				query += `AND created <= $2 LIMIT NULLIF($3, 0)`
			}
		}
		rows, err = transaction.Query(context.Background(), query, userID, params.Since, params.Limit)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	userTransactions := models.Transactions{}
	var receiver sql.NullInt64
	for rows.Next() {
		var userTransaction models.Transaction
		if err = rows.Scan(&userTransaction.OperationType, &receiver, &userTransaction.Amount,
			&userTransaction.Created); err != nil {
			return nil, err
		}

		if receiver.Valid {
			userTransaction.ReceiverID = receiver.Int64
		} else {
			userTransaction.ReceiverID = 0
		}
		userTransactions = append(userTransactions, &userTransaction)
	}

	return userTransactions, nil
}

func (s *Storage) DoesUserExist(userID int64) (bool, error) {
	transaction, err := s.db.Begin(context.Background())
	defer func() {
		if err != nil {
			_ = transaction.Rollback(context.Background())
		} else {
			_ = transaction.Commit(context.Background())
		}
	}()

	if err = transaction.QueryRow(context.Background(), queryGetUserID, userID).Scan(&userID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return false, err
		}
		err = nil
		return false, nil
	}
	return true, nil
}
