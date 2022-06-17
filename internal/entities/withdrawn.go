package entities

import (
	"github.com/jackc/pgtype"
)

type Withdraw struct {
	ID          int
	UserID      int
	Order       string
	Sum         int
	ProcessedAt pgtype.Timestamp
}
