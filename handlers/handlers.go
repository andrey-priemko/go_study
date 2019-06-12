package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type VideoListItem struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Duration int `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

type Video struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Duration int `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	Url string `json:"url"`
}

const dirPath = "content"

const videoContentType = "video/mp4"
const videoFileName = "index.mp4"
const thumbFileName = "screen.jpg"

func Router(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/list", getListFromDb(db)).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", getVideoFromDb(db)).Methods(http.MethodGet)
	s.HandleFunc("/video", uploadVideoIntoDb(db)).Methods(http.MethodPost)

	return logMiddleware(r)
}

func createFile(id string) (*os.File, error) {
	if err := os.Mkdir(dirPath, os.ModeDir); err != nil && !os.IsExist(err) {
		return nil, err
	}

	videoDirPath := filepath.Join(dirPath, id)
	if err := os.Mkdir(videoDirPath, os.ModeDir); err != nil && !os.IsExist(err) {
		return nil, err
	}

	filePath := filepath.Join(videoDirPath, videoFileName)
	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}

func getListFromDb(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT video_key, title, duration, thumbnail_url FROM video")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		videos := make([]VideoListItem, 0)
		for rows.Next() {
			var videoListItem VideoListItem
			err := rows.Scan(
				&videoListItem.Id,
				&videoListItem.Name,
				&videoListItem.Duration,
				&videoListItem.Thumbnail,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			videos = append(videos, videoListItem)
		}

		writeResponseData(w, videos)
	}
}

func getVideoFromDb(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["ID"]

		var video Video
		err := db.QueryRow("SELECT video_key, title, duration, thumbnail_url, url FROM video WHERE video_key = ?", id).Scan(
			&video.Id,
			&video.Name,
			&video.Duration,
			&video.Thumbnail,
			&video.Url,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeResponseData(w, video)
	}
}

func uploadVideoIntoDb(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileReader, header, err := r.FormFile("file[]")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contentType := header.Header.Get("Content-Type")
		if contentType != videoContentType {
			http.Error(w, "Unexpected content type", http.StatusBadRequest)
			return
		}

		fileName := header.Filename

		id, err := uuid.NewUUID()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		idStr := id.String()

		file, err := createFile(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close() // случайно заигнорил тип ошибки unhandled error

		_, err = io.Copy(file, fileReader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		videoPath := filepath.Join(dirPath, idStr, videoFileName)
		thumbPath := filepath.Join(dirPath, idStr, thumbFileName)

		db.QueryRow(
			"INSERT INTO video SET video_key = ?, title = ?, duration = ?, thumbnail_url = ?, url = ?, status = ?",
			idStr,
			fileName,
			127,
			thumbPath,
			videoPath,
			3,
		)

		w.WriteHeader(http.StatusOK)
	}
}

func writeResponseData(w http.ResponseWriter, data interface{})  {
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err = io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}
