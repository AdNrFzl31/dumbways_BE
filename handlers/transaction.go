package handlers

import (
	dto "dumbsound/dto/result"
	transactiondto "dumbsound/dto/transaction"
	"dumbsound/models"
	"dumbsound/repositories"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contetnt-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))
	startDate := time.Now()
	dueDate := startDate.AddDate(0, 1, 0)
	request := transactiondto.TransactionRequest{
		StartDate: startDate,
		DueDate:   dueDate,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	setTransID := time.Now().Unix()

	transaction := models.Transaction{
		ID:        int(setTransID),
		UserID:    userId,
		StartDate: request.StartDate,
		DueDate:   request.DueDate,
		Status:    "pending",
		Price:     48999,
	}

	data, err := h.TransactionRepository.CreateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Non midtrans below

	// data, _ = h.TransactionRepository.GetTransactionID(data.ID)

	// format := "2006-01-02"

	// transactionResponse := transactiondto.TransactionResponse{
	// 	StartDate: data.StartDate.Format(format),
	// 	DueDate:   data.DueDate.Format(format),
	// 	User:      data.User,
	// 	Status:    data.Status,
	// 	Price:     data.Price,
	// }

	// w.WriteHeader(http.StatusOK)
	// response := dto.SuccessResult{Status: "success", Data: transactionResponse}
	// json.NewEncoder(w).Encode(response)

	// Midtrans
	DataSnap, _ := h.TransactionRepository.GetTransactionID(data.ID)

	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(int(DataSnap.ID)),
			GrossAmt: int64(DataSnap.Price),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: DataSnap.User.Fullname,
			Email: DataSnap.User.Email,
		},
	}

	// Run midtrans Snap

	snapResp, _ := s.CreateTransaction(req)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "success", Data: snapResp}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderID := notificationPayload["order_id"].(string)

	transaction, _ := h.TransactionRepository.GetTransactionMidtrans(orderID)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			h.TransactionRepository.UpdateTransactionStatus("pending", int(transaction.ID))
			h.TransactionRepository.UpdateUserSubscribe("false", userId)
		} else if fraudStatus == "accept" {
			h.TransactionRepository.UpdateTransactionStatus("success", int(transaction.ID))
			h.TransactionRepository.UpdateUserSubscribe("true", userId)
		}

	} else if transactionStatus == "settlement" {
		h.TransactionRepository.UpdateTransactionStatus("success", int(transaction.ID))
		h.TransactionRepository.UpdateUserSubscribe("true", userId)
	} else if transactionStatus == "deny" {
		h.TransactionRepository.UpdateTransactionStatus("failed", int(transaction.ID))
		h.TransactionRepository.UpdateUserSubscribe("false", userId)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		h.TransactionRepository.UpdateTransactionStatus("failed", int(transaction.ID))
		h.TransactionRepository.UpdateUserSubscribe("false", userId)
	} else if transactionStatus == "pending" {
		h.TransactionRepository.UpdateTransactionStatus("pending", int(transaction.ID))
		h.TransactionRepository.UpdateUserSubscribe("false", userId)
	}
	w.WriteHeader(http.StatusOK)
}
