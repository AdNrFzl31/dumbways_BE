package middleware

import (
	"context"
	dto "dumbsound/dto/result"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func UploadMusic(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("music")

		if err != nil && r.Method == "PATCH" {
			ctx := context.WithValue(r.Context(), "dataMusic", "false")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode("Error Retrieving the File Music")
			return
		}
		defer file.Close()

		// // setup file type filtering
		// buff := make([]byte, 512)
		// _, err = file.Read(buff)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
		// 	json.NewEncoder(w).Encode(response)
		// 	return
		// }

		// filetype := http.DetectContentType(buff)
		// if filetype != "audio/mp3" {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	response := dto.ErrorResult{Status: "Failed", Message: "The provided file format is not allowed. Please upload a MP3 music"}
		// 	json.NewEncoder(w).Encode(response)
		// 	return
		// }

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Status: "Server Error", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		const MAX_UPLOAD_SIZE = 50 << 20
		r.ParseMultipartForm(MAX_UPLOAD_SIZE)
		if r.ContentLength > MAX_UPLOAD_SIZE {
			w.WriteHeader(http.StatusBadRequest)
			response := Result{Status: "Failed", Message: "Max size in 50mb"}
			json.NewEncoder(w).Encode(response)
			return
		}
		tempFile, err := ioutil.TempFile("uploads", "music-*.mp3")
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		tempFile.Write(fileBytes)

		data := tempFile.Name()
		// filepdf := data[8:]

		ctx := context.WithValue(r.Context(), "dataMusic", data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
