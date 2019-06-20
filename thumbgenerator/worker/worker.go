package worker

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go_study/thumbgenerator/model"
	"go_study/thumbgenerator/processor"
)

func Worker(tasksChan <-chan *model.Task, db *sql.DB, name int) {
	log.Printf("start worker %v\n", name)
	for task := range tasksChan {
		log.Printf("start processing video with id %v on worker %v\n", task.Id, name)
		processor.ProcessTask(task, db)
		log.Printf("end processing video with id %v on worker %v\n", task.Id, name)
	}
	log.Printf("stop worker %v\n", name)
}