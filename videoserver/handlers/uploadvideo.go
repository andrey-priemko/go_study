package handlers

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go_study/videoserver/provider"
	"go_study/videoserver/storage"
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

		id, err := uuid.NewUUID()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		videoId := id.String()

		err = storage.CopyFile(fileReader, videoId)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = provider.UploadVideo(
			videoId,
			header.Filename,
			filepath.Join(storage.ContentDir, videoId, storage.VideoFileName),
		)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}