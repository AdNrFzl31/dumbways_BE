package repositories

import (
	"dumbsound/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactions() ([]models.Transaction, error)

	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	GetTransactionID(ID int) (models.Transaction, error)
	UpdateTransactionStatus(status string, ID int) error
	UpdateUserSubscribe(subscribe string, ID int) error
	GetTransactionMidtrans(ID int) (models.Transaction, error)
	CancelTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(transaction models.Transaction) (models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}
func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Find(&transaction).Error
	return transaction, err
}

func (r *repository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("User").Create(&transaction).Error

	return transaction, err
}

func (r *repository) GetTransactionID(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Find(&transaction, ID).Error

	return transaction, err
}

func (r *repository) UpdateTransactionStatus(status string, ID int) error {
	var transaction models.Transaction
	r.db.Preload("User").First(&transaction.ID)

	transaction.Status = status
	err := r.db.Save(&transaction).Error

	return err
}

func (r *repository) UpdateUserSubscribe(subscribe string, ID int) error {
	var user models.User
	r.db.Preload("User").First(&user, ID)

	user.Subscribe = subscribe
	err := r.db.Save(&user).Error

	return err
}

func (r *repository) GetTransactionMidtrans(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, ID).Error
	return transaction, err
}

func (r *repository) CancelTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Save(&transaction).Error
	return transaction, err
}

func (r *repository) UpdateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Save(&transaction).Error
	return transaction, err
}
