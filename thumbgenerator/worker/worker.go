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
		log.Printf("start handle task %v on worker %v\n", task.VideoKey, name)
		processor.ProcessTask(task, db)
		log.Printf("end handle task %v on worker %v\n", task.VideoKey, name)
	}
	log.Printf("stop worker %v\n", name)
}