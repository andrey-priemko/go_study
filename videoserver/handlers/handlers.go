package handlers

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go_study/videoserver/provider"
	"net/http"
)

func Router(dp provider.DataProvider) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/list", getVideoList(dp)).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", getVideo(dp)).Methods(http.MethodGet)
	s.HandleFunc("/video", uploadVideo(dp)).Methods(http.MethodPost)

	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}