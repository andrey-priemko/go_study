package mysql

import (
	"go_study/videoserver/model"
)

func (c *Connector) GetVideo(id string) (*model.Video, error) {
	var video model.Video
	err := c.DB.QueryRow("SELECT video_key, title, duration, thumbnail_url, url FROM video WHERE video_key=?", id).Scan(
		&video.Id,
		&video.Name,
		&video.Duration,
		&video.Thumbnail,
		&video.Url,
	)
	return &video, err
}