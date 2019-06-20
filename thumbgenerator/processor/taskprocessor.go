package processor

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go_study/thumbgenerator/model"
	"go_study/videoserver/database"
	"go_study/videoserver/storage"
	"os/exec"
	"path/filepath"
	"strconv"
)

func ProcessTask(task *model.Task, db *sql.DB) {
	thumbUrl := filepath.Join(filepath.Dir(task.Url), storage.ThumbFileName)
	out, err := exec.Command("D:\\projects\\go_dev\\dev\\src\\go_study\\bin\\VideoProcessor.exe", task.Url, thumbUrl).Output()
	if err != nil {
		log.Error(err.Error())
		return
	}

	duration, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = database.ExecTransaction(
		db,
		"UPDATE video SET status=?, duration=?, thumbnail_url=? WHERE id=?",
		model.Processed,
		duration,
		thumbUrl,
		task.Id,
	)
	if err != nil {
		log.Error(err.Error())
	}
}
