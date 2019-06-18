package handlers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"go_study/videoserver/model"
	"net/http"
)

func getVideo(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["ID"]

		var video model.Video
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
