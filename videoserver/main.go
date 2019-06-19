package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go_study/videoserver/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	db, err := sql.Open("mysql", `root:1234@/go_dev`)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()

	killSignalChan := getKillSignalChan()

	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
	srv := startServer(serverUrl, db)

	waitForKillSignal(killSignalChan)
	srv.Shutdown(context.Background())
}

func startServer(serverUrl string, db *sql.DB) *http.Server {
	router := handlers.Router(db)
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Error(srv.ListenAndServe())
	}()

	return srv
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}