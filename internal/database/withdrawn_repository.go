package database

import (
	"context"
	"fmt"
	"github.com/AtCliffUnderline/golang-diploma/internal/entities"
	"time"
)

type WithdrawnRepository interface {
	Add(userID int, order string, sum int) (bool, error)
	GetByOrderID(orderID int) (entities.Withdraw, error)
	GetByUserID(userID int) ([]entities.Withdraw, error)
	GetUserWithdrawnSum(userID int) (int, error)
}

type WithdrawnStorage struct {
	Storage Storage
}

func (wr WithdrawnStorage) getTableName() string {
	return "withdrawals"
}

func (wr WithdrawnStorage) Add(userID int, order string, sum int) (bool, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, order_number, sum, processed_at) VALUES (%d, %s, %d, '%s');", wr.getTableName(), userID, order, sum, time.Now().Format(time.RFC3339))
	_, err := wr.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (wr WithdrawnStorage) GetByOrderID(orderID int) (entities.Withdraw, error) {
	var res entities.Withdraw
	query := fmt.Sprintf("SELECT id, user_id, order_number, sum, processed_at FROM %s WHERE order_id = %d;", wr.getTableName(), orderID)
	err := wr.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res.ID, &res.UserID, &res.Order, &res.Sum, &res.ProcessedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (wr WithdrawnStorage) GetByUserID(userID int) ([]entities.Withdraw, error) {
	res := make([]entities.Withdraw, 0)
	query := fmt.Sprintf("SELECT id, user_id, order_number, sum, processed_at FROM %s WHERE user_id = %d ORDER BY w.processed_at ASC;", wr.getTableName(), userID)
	rows, err := wr.Storage.DBConn.Query(context.Background(), query)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var r entities.Withdraw
		err := rows.Scan(&r.ID, &r.UserID, &r.Order, &r.Sum, &r.ProcessedAt)
		if err != nil {
			return nil, nil
		}
		res = append(res, r)
	}

	return res, nil
}

func (wr WithdrawnStorage) GetUserWithdrawnSum(userID int) (int, error) {
	var res int
	query := fmt.Sprintf("SELECT COALESCE(SUM(sum), 0) FROM %s WHERE user_id = %d;", wr.getTableName(), userID)
	err := wr.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res)
	if err != nil {
		return 0, err
	}

	return res, nil
}
