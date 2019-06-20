package mysql

import (
	"go_study/videoserver/database"
)

func (c *Connector) UploadVideo(id string, fileName string, url string) error {
	return database.ExecTransaction(
		c.DB,
		"INSERT INTO video SET video_key=?, title=?, url=?",
		id,
		fileName,
		url,
	)
}
