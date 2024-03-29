package provider

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go_study/thumbgenerator/generator"
	"go_study/thumbgenerator/model"
	"time"
)

func TaskProvider(stopChan chan struct{}, db *sql.DB) <-chan *model.Task {
	tasksChan := make(chan *model.Task)
	go func() {
		for {
			select {
			case <-stopChan:
				close(tasksChan)
				return
			default:
			}
			if task := generator.GenerateTask(db); task != nil {
				log.Printf("got the task %v\n", task)
				tasksChan <- task
			} else {
				log.Info("no task for processing, start waiting")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return tasksChan
}

func RunTaskProvider(stopChan chan struct{}, db *sql.DB) <-chan *model.Task {
	resultChan := make(chan *model.Task)
	stopTaskProviderChan := make(chan struct{})
	taskProviderChan := TaskProvider(stopTaskProviderChan, db)
	onStop := func () {
		stopTaskProviderChan <- struct{}{}
		close(resultChan)
	}
	go func() {
		for {
			select {
			case <-stopChan:
				onStop()
				return
			case task := <-taskProviderChan:
				select {
				case <-stopChan:
					onStop()
					return
				case resultChan <- task:
				}
			}
		}
	}()
	return resultChan
}
