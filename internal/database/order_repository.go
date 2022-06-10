package database

import (
	"context"
	"fmt"
	"github.com/AtCliffUnderline/golang-diploma/internal/entities"
	"time"
)

type OrderRepository interface {
	Add(userID int, number string, status string, accrual int) (bool, error)
	GetByNumber(number string) (entities.Order, error)
	GetByUserID(userID int) ([]entities.Order, error)
	Update(number string, status string, accrual int) (bool, error)
	GetAllNonFinalOrders() ([]entities.Order, error)
}

type OrderStorage struct {
	Storage Storage
}

func (or OrderStorage) GetTableName() string {
	return "orders"
}

func (or OrderStorage) Add(userID int, number string, status string, accrual int) (bool, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, number, status, accrual, uploaded_at) VALUES (%d, '%s', '%s', %d, '%s');", or.GetTableName(), userID, number, status, accrual, time.Now().Format(time.RFC3339))
	_, err := or.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (or OrderStorage) GetByNumber(number string) (entities.Order, error) {
	var res entities.Order
	query := fmt.Sprintf("SELECT id, user_id, number, status, accrual, uploaded_at FROM %s WHERE number = '%s';", or.GetTableName(), number)
	err := or.Storage.DBConn.QueryRow(context.Background(), query).Scan(&res.ID, &res.UserID, &res.Number, &res.Status, &res.Accrual, &res.UploadedAt)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (or OrderStorage) GetByUserID(userID int) ([]entities.Order, error) {
	res := make([]entities.Order, 0)
	query := fmt.Sprintf("SELECT id, user_id, number, status, accrual, uploaded_at FROM %s WHERE user_id = %d ORDER BY uploaded_at ASC;", or.GetTableName(), userID)
	rows, err := or.Storage.DBConn.Query(context.Background(), query)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var r entities.Order
		err := rows.Scan(&r.ID, &r.UserID, &r.Number, &r.Status, &r.Accrual, &r.UploadedAt)
		if err != nil {
			return nil, nil
		}
		res = append(res, r)
	}

	return res, nil
}

func (or OrderStorage) Update(number string, status string, accrual int) (bool, error) {
	query := fmt.Sprintf("UPDATE %s SET status = '%s', accrual = %d WHERE number = '%s';", or.GetTableName(), status, accrual, number)
	_, err := or.Storage.DBConn.Exec(context.Background(), query)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (or OrderStorage) GetAllNonFinalOrders() ([]entities.Order, error) {
	res := make([]entities.Order, 0)
	query := fmt.Sprintf("SELECT id, user_id, number, status, accrual, uploaded_at FROM %s WHERE status NOT IN ('PROCESSED', 'INVALID');", or.GetTableName())
	rows, err := or.Storage.DBConn.Query(context.Background(), query)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var r entities.Order
		err := rows.Scan(&r.ID, &r.UserID, &r.Number, &r.Status, &r.Accrual, &r.UploadedAt)
		if err != nil {
			return nil, nil
		}
		res = append(res, r)
	}

	return res, nil
}
