package entities

import (
	"github.com/jackc/pgtype"
)

type Order struct {
	ID         int
	UserID     int
	Number     string
	Status     string
	Accrual    int
	UploadedAt pgtype.Timestamp
}
