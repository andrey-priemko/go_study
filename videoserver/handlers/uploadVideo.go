package handlers

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go_study/videoserver/storage"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
)

const DIR_PATH = "content"

const VIDEO_CONTENT_TYPE = "video/mp4"
const VIDEO_FILE_NAME = "index.mp4"
const THUMB_FILE_NAME = "screen.jpg"

func uploadVideo(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileReader, header, err := r.FormFile("file[]")
		if err != nil {
			log.Fatal(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contentType := header.Header.Get("Content-Type")
		if contentType != VIDEO_CONTENT_TYPE {
			http.Error(w, "Unexpected content type", http.StatusBadRequest)
			return
		}

		file, id, err := storage.CreateFile(VIDEO_FILE_NAME)
		if err != nil {
			log.Fatal(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, fileReader)
		if err != nil {
			log.Fatal(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		videoPath := filepath.Join(DIR_PATH, id, VIDEO_FILE_NAME)
		thumbPath := filepath.Join(DIR_PATH, id, THUMB_FILE_NAME)

		out, err := exec.Command("D:\\projects\\go_dev\\dev\\src\\go_study\\VideoProcessor.exe", videoPath, thumbPath).Output() //todo
		if err != nil {
			log.Fatal(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		duration, err := strconv.ParseFloat(string(out), 64)
		if err != nil {
			log.Fatal(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		defer tx.Rollback()

		rows, err := tx.Query(
			"INSERT INTO video SET video_key = ?, title = ?, duration = ?, thumbnail_url = ?, url = ?, status = ?",
			id,
			header.Filename,
			duration,
			thumbPath,
			videoPath,
			3,
		)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		defer rows.Close()

		err = tx.Commit()
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}