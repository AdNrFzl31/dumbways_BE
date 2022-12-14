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
	"github.com/gorilla/mux"
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

	transID := time.Now().Unix()

	transaction := models.Transaction{
		ID:        int(transID),
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

	DataSnap, err := h.TransactionRepository.GetTransactionID(data.ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

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

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderID := notificationPayload["order_id"].(string)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			h.TransactionRepository.UpdateTransactionStatus("pending", orderID)

		} else if fraudStatus == "accept" {
			h.TransactionRepository.UpdateTransactionStatus("success", orderID)
		}

	} else if transactionStatus == "settlement" {
		h.TransactionRepository.UpdateTransactionStatus("success", orderID)

	} else if transactionStatus == "deny" {
		h.TransactionRepository.UpdateTransactionStatus("failed", orderID)

	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		h.TransactionRepository.UpdateTransactionStatus("failed", orderID)

	} else if transactionStatus == "pending" {
		h.TransactionRepository.UpdateTransactionStatus("pending", orderID)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handlerTransaction) FindTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transactions, err := h.TransactionRepository.FindTransaction()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: transactions}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) CancelTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	transaction, err := h.TransactionRepository.GetTransactionID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transaction.Status = "Cancel"
	data, err := h.TransactionRepository.CancelTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) AcceptTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	transaction, err := h.TransactionRepository.GetTransactionID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: "Cek id Transaction => " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// fmt.Println(transaction.ID)
	// fmt.Println(transaction.Status)

	transaction.Status = "Success"
	data, err := h.TransactionRepository.UpdateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: data}
	json.NewEncoder(w).Encode(response)
}
