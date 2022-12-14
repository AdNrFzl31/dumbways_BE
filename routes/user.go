package routes

import (
	"dumbsound/handlers"
	"dumbsound/pkg/middleware"
	"dumbsound/pkg/mysql"
	"dumbsound/repositories"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	UserRepository := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerUser(UserRepository)

	r.HandleFunc("/users", h.FindUsers).Methods("GET")
	r.HandleFunc("/user/{id}", h.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", middleware.Auth(middleware.UploadProfile(h.UpdateUser))).Methods("PATCH")
	r.HandleFunc("/user/{id}", middleware.Auth(h.DeleteUser)).Methods("DELETE")
}
