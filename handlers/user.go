package handlers

import (
	dto "dumbsound/dto/result"
	usersdto "dumbsound/dto/users"
	"dumbsound/models"
	"dumbsound/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerUser struct {
	userRepository repositories.UserRepository
}

func HandlerUser(UserRepository repositories.UserRepository) *handlerUser {
	return &handlerUser{UserRepository}
}

func (h *handlerUser) FindUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.userRepository.FindUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: users}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.userRepository.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: user}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userStatus := userInfo["status"]
	userId := int(userInfo["id"].(float64))

	if userId != id && userStatus != "admin" {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: "you're not admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.userRepository.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.userRepository.DeleteUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: convertResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func convertResponse(u models.User) usersdto.UserResponse {
	return usersdto.UserResponse{
		ID:       u.ID,
		Fullname: u.Fullname,
		Email:    u.Email,
		// Image:    u.Image,
	}
}
