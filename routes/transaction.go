package routes

import (
	"dumbsound/handlers"
	"dumbsound/pkg/middleware"
	"dumbsound/pkg/mysql"
	"dumbsound/repositories"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router) {
	transactionRepository := repositories.RepositoryTransaction(mysql.DB)
	h := handlers.HandlerTransaction(transactionRepository)

	r.HandleFunc("/transall", h.FindTransactions).Methods("GET")
	r.HandleFunc("/transaction", middleware.Auth(h.CreateTransaction)).Methods("POST")
	r.HandleFunc("/canceltrans/{id}", middleware.Auth(h.CancelTransaction)).Methods("PATCH")
	r.HandleFunc("/accepttrans/{id}", middleware.Auth(h.AcceptTransaction)).Methods("PATCH")
	r.HandleFunc("/notification", middleware.Auth(h.Notification)).Methods("POST")
}
