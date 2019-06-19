package handlers

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go_study/videoserver/storage"
	"io"
	"net/http"
	"path/filepath"
)

func uploadVideo(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileReader, header, err := r.FormFile("file[]")
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contentType := header.Header.Get("Content-Type")
		if contentType != storage.VIDEO_CONTENT_TYPE {
			log.Error("Unexpected content type", contentType)
			http.Error(w, "Unexpected content type", http.StatusBadRequest)
			return
		}

		file, id, err := storage.CreateFile(storage.VIDEO_FILE_NAME)
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

		tx, err := db.Begin()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer tx.Rollback()

		rows, err := tx.Query(
			"INSERT INTO video SET video_key = ?, title = ?, url = ?",
			id,
			header.Filename,
			filepath.Join(storage.CONTENT_DIR, id, storage.VIDEO_FILE_NAME),
		)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()

		err = tx.Commit()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}