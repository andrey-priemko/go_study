package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go_study/thumbgenerator/model"
	"go_study/videoserver/database"
	"go_study/videoserver/storage"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const WORKERS_COUNT = 3

type Task struct {
	VideoKey string
	Url string
}

func GenerateTask(db *sql.DB) *Task {
	var task Task
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

func TaskProvider(stopChan chan struct{}, db *sql.DB) <-chan *Task {
	tasksChan := make(chan *Task)
	go func() {
		for {
			select {
			case <-stopChan:
				close(tasksChan)
				return
			default:
			}
			if task := GenerateTask(db); task != nil {
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

func RunTaskProvider(stopChan chan struct{}, db *sql.DB) <-chan *Task {
	resultChan := make(chan *Task)
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

func Worker(tasksChan <-chan *Task, db *sql.DB, name int) {
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

func RunWorkerPool(stopChan chan struct{}, db *sql.DB) *sync.WaitGroup {
	var wg sync.WaitGroup
	tasksChan := RunTaskProvider(stopChan, db)
	for i := 0; i < WORKERS_COUNT; i++ {
		go func(i int) {
			wg.Add(1)
			Worker(tasksChan, db, i)
			wg.Done()
		}(i)
	}
	return &wg
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("thumbgenerator.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	db, err := sql.Open("mysql", `root:1234@/go_dev`)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()

	rand.Seed(time.Now().Unix())
	stopChan := make(chan struct{})

	killChan := getKillSignalChan()
	wg := RunWorkerPool(stopChan, db)

	waitForKillSignal(killChan)
	stopChan <- struct{}{}
	wg.Wait()
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}