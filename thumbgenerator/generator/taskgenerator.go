package generator

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go_study/thumbgenerator/model"
	"go_study/videoserver/database"
)

func GenerateTask(db *sql.DB) *model.Task {
	var task model.Task
	err := db.QueryRow("SELECT video_key, url FROM video WHERE status=?", model.NotProcessed).Scan(
		&task.VideoKey,
		&task.Url,
	)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	err = database.ExecTransaction(
		db,
		"UPDATE video SET status=? WHERE video_key=?", model.Processing,
		task.VideoKey,
	)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return &task
}