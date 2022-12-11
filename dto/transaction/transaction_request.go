package transactiondto

import "time"

type TransactionRequest struct {
	StartDate time.Time `json:"startDate" `
	DueDate   time.Time `json:"dueDate"`
	UserID    int       `json:"userId"`
}
