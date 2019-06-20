package handlers

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go_study/videoserver/provider"
	"net/http"
)

func getVideo(dp provider.DataProvider) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["ID"]

		video, err := dp.GetVideo(id)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeResponseData(w, video)
	}
}
