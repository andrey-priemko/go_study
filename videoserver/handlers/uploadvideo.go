package handlers

import (
	log "github.com/sirupsen/logrus"
	"go_study/videoserver/provider"
	"go_study/videoserver/storage"
	"io"
	"net/http"
	"path/filepath"
)

func uploadVideo(provider provider.DataProvider) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileReader, header, err := r.FormFile("file[]")
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contentType := header.Header.Get("Content-Type")
		if contentType != storage.VideoContentType {
			log.Error("Unexpected content type", contentType)
			http.Error(w, "Unexpected content type", http.StatusBadRequest)
			return
		}

		file, id, err := storage.CreateFile(storage.VideoFileName)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, fileReader)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = provider.UploadVideo(
			id,
			header.Filename,
			filepath.Join(storage.ContentDir, id, storage.VideoFileName),
		)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}