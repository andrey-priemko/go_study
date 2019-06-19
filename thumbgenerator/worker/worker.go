package worker

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

func Worker(tasksChan <-chan *model.Task, db *sql.DB, name int) {
	log.Printf("start worker %v\n", name)
	for task := range tasksChan {
		log.Printf("start handle task %v on worker %v\n", task.VideoKey, name)

		thumbUrl := filepath.Join("content", task.VideoKey, storage.ThumbFileName) //todo
		out, err := exec.Command("D:\\projects\\go_dev\\dev\\src\\go_study\\bin\\VideoProcessor.exe", task.Url, thumbUrl).Output()
		if err != nil {
			log.Error(err.Error())
			continue
		}

		duration, err := strconv.ParseFloat(string(out), 64)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		log.Printf("duration", duration)

		err = database.ExecTransaction(
			db,
			"UPDATE video SET status=?, duration=?, thumbnail_url=? WHERE video_key=?",
			model.Processed,
			duration,
			thumbUrl,
			task.VideoKey,
		)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		log.Printf("end handle task %v on worker %v\n", task.VideoKey, name)
	}
	log.Printf("stop worker %v\n", name)
}