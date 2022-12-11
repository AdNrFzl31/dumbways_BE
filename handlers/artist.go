package handlers

import (
	artistdto "dumbsound/dto/artist"
	dto "dumbsound/dto/result"
	"dumbsound/models"
	"dumbsound/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerArtist struct {
	ArtistRepository repositories.ArtistRepository
}

func HandlerArtist(ArtistRepository repositories.ArtistRepository) *handlerArtist {
	return &handlerArtist{ArtistRepository}
}

func (h *handlerArtist) FindArtists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	artists, err := h.ArtistRepository.FindArtists()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: artists}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerArtist) GetArtist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var artist models.Artist
	artist, err := h.ArtistRepository.GetArtist(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: artist}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerArtist) CreateArtist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userStatus := userInfo["status"]

	if userStatus != "admin" {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: "You're Not Admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var request models.Artist

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	artist := models.Artist{
		Name:   request.Name,
		Old:    request.Old,
		Artist: request.Artist,
		Career: request.Career,
	}

	data, err := h.ArtistRepository.CreateArtist(artist)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "success", Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerArtist) UpdateArtist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request models.Artist
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	artist, err := h.ArtistRepository.GetArtist(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// artist := models.Artist{}

	if request.Name != "" {
		artist.Name = request.Name
	}

	if request.Old != 0 {
		artist.Old = request.Old
	}

	if request.Artist != "" {
		artist.Artist = request.Artist
	}

	if request.Career != "" {
		artist.Career = request.Career
	}

	data, err := h.ArtistRepository.UpdateArtist(artist)
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

func (h *handlerArtist) DeleteArtist(w http.ResponseWriter, r *http.Request) {
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

	artists, err := h.ArtistRepository.GetArtist(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	deleteArtist, err := h.ArtistRepository.DeleteArtist(artists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: convertResponseArtist(deleteArtist)}
	json.NewEncoder(w).Encode(response)
}

func convertResponseArtist(u models.Artist) artistdto.ArtistResponseDelete {
	return artistdto.ArtistResponseDelete{
		ID: u.ID,
	}
}
