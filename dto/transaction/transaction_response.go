package transactiondto

import "dumbsound/models"

type TransactionResponse struct {
	StartDate string      `json:"startDate" `
	DueDate   string      `json:"dueDate"`
	User      models.User `json:"user"`
	Status    string      `json:"status"`
	Price     int         `json:"price"`
}
