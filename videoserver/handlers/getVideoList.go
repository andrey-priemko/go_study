package handlers

import (
	"database/sql"
	"go_study/videoserver/model"
	"net/http"
)

func getVideoList(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT video_key, title, duration, thumbnail_url FROM video")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		videos := make([]model.VideoListItem, 0)
		for rows.Next() {
			var videoListItem model.VideoListItem
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