package mysql

import (
	"go_study/videoserver/model"
)

func (c *Connector) GetVideoList() ([]model.VideoListItem, error) {
	rows, err := c.DB.Query("SELECT video_key, title, duration, thumbnail_url FROM video")
	if err != nil {
		return nil, err
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
			return nil, err
		}
		videos = append(videos, videoListItem)
	}

	return videos, nil
}