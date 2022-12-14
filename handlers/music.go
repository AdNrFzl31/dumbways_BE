package handlers

import (
	"context"
	musicdto "dumbsound/dto/music"
	dto "dumbsound/dto/result"
	"dumbsound/models"
	"dumbsound/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerMusic struct {
	musicRepository repositories.MusicRepository
}

func HandlerMusic(MusicRepository repositories.MusicRepository) *handlerMusic {
	return &handlerMusic{MusicRepository}
}

func (h *handlerMusic) FindMusics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	musics, err := h.musicRepository.FindMusic()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: musics}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerMusic) GetMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var music models.Music
	music, err := h.musicRepository.GetMusic(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: music}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerMusic) CreateMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// ambil data user id dari token yang sudah di decode
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userStatus := userInfo["status"]
	if userStatus != "admin" {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: "You're Not Admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get dataFile from midleware and store to filename variable
	dataContextFile := r.Context().Value("dataFile")
	filepath := dataContextFile.(string)

	dataContextMusic := r.Context().Value("dataMusic")
	musicpath := dataContextMusic.(string)

	year, _ := strconv.Atoi(r.FormValue("year"))
	artistId, _ := strconv.Atoi(r.FormValue("artistId"))
	fmt.Println(artistId)
	request := musicdto.MusicRequest{
		Title:    r.FormValue("title"),
		Year:     year,
		ArtistId: artistId,
		Tumbnail: filepath,
		Music:    musicpath,
	}
	fmt.Println(request)

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Declare Context Background, Cloud Name, API Key, API Secret ...
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	respFile, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "dumbsound"})

	if err != nil {
		fmt.Println(err.Error())
	}

	respMusic, err := cld.Upload.Upload(ctx, musicpath, uploader.UploadParams{Folder: "dumbsound"})

	if err != nil {
		fmt.Println(err.Error())
	}

	music := models.Music{
		Title:    request.Title,
		Year:     request.Year,
		ArtistId: request.ArtistId,
		Tumbnail: respFile.SecureURL,
		Music:    respMusic.SecureURL,
		// UserID:      userId,
	}

	music, err = h.musicRepository.CreateMusic(music)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	music, _ = h.musicRepository.GetMusic(music.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: music}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerMusic) UpdateMusic(w http.ResponseWriter, r *http.Request) {
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

	dataContextFile := r.Context().Value("dataFile")
	filepath := dataContextFile.(string)

	dataContextMusic := r.Context().Value("dataMusic")
	musicpath := dataContextMusic.(string)

	year, _ := strconv.Atoi(r.FormValue("year"))
	artistId, _ := strconv.Atoi(r.FormValue("artistId"))
	request := musicdto.MusicRequest{
		Title:    r.FormValue("title"),
		Year:     year,
		ArtistId: artistId,
		Tumbnail: filepath,
		Music:    musicpath,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Declare Context Background, Cloud Name, API Key, API Secret ...
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	respFile, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "dumbsound"})

	if err != nil {
		fmt.Println(err.Error())
	}

	respMusic, err := cld.Upload.Upload(ctx, musicpath, uploader.UploadParams{Folder: "dumbsound"})

	if err != nil {
		fmt.Println(err.Error())
	}

	music, err := h.musicRepository.GetMusic(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// music := models.Music{}

	if request.Title != "" {
		music.Title = request.Title
	}

	if request.Year != 0 {
		music.Year = request.Year
	}

	if request.Tumbnail != "" {
		music.Tumbnail = respFile.SecureURL
	}

	if request.ArtistId != 0 {
		music.ArtistId = request.ArtistId
	}

	if request.Music != "" {
		music.Music = respMusic.SecureURL
	}

	data, err := h.musicRepository.UpdateMusic(music)
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

func (h *handlerMusic) DeleteMusic(w http.ResponseWriter, r *http.Request) {
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

	musics, err := h.musicRepository.GetMusic(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: "Failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	deleteMusic, err := h.musicRepository.DeleteMusic(musics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: "Success", Data: convertResponseMusic(deleteMusic)}
	json.NewEncoder(w).Encode(response)
}

func convertResponseMusic(u models.Music) musicdto.MusicResponseDelete {
	return musicdto.MusicResponseDelete{
		ID: u.ID,
	}
}
