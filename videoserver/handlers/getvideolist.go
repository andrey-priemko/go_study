package handlers

import (
	log "github.com/sirupsen/logrus"
	"go_study/videoserver/provider"
	"net/http"
)

func getVideoList(dp provider.DataProvider) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		videos, err := dp.GetVideoList()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeResponseData(w, videos)
	}
}