package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type VideoListItem struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Duration int `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

type Video struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Duration int `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	Url string `json:"url"`
}

const dirPath = "content"

const videoContentType = "video/mp4"

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/list", getList).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", getVideo).Methods(http.MethodGet)
	s.HandleFunc("/video", uploadVideo).Methods(http.MethodPost)

	return logMiddleware(r)
}

func getList(w http.ResponseWriter, _ *http.Request) {
	videoList := []VideoListItem{
		{
			"d290f1ee-6c54-4b01-90e6-d701748f0851",
			"Black Retrospetive Woman",
			15,
			"/content/d290f1ee-6c54-4b01-90e6-d701748f0851/screen.jpg",
		},
		{
			"sldjfl34-dfgj-523k-jk34-5jk3j45klj34",
			"Go Rally TEASER-HD",
			41,
			"/content/sldjfl34-dfgj-523k-jk34-5jk3j45klj34/screen.jpg",
		},
		{
			"hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345",
			"Танцор",
			92,
			"/content/hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345/screen.jpg",
		},
	}

	b, err := json.Marshal(videoList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err = io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func getVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]

	videoList := []Video{
		{
			"d290f1ee-6c54-4b01-90e6-d701748f0851",
			"Black Retrospetive Woman",
			15,
			"/content/d290f1ee-6c54-4b01-90e6-d701748f0851/screen.jpg",
			"/content/d290f1ee-6c54-4b01-90e6-d701748f0851/index.mp4",
		},
		{
			"sldjfl34-dfgj-523k-jk34-5jk3j45klj34",
			"Go Rally TEASER-HD",
			41,
			"/content/sldjfl34-dfgj-523k-jk34-5jk3j45klj34/screen.jpg",
			"/content/sldjfl34-dfgj-523k-jk34-5jk3j45klj34/index.mp4",
		},
		{
			"hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345",
			"Танцор",
			92,
			"/content/hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345/screen.jpg",
			"/content/hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345/index.mp4",
		},
	}
	var video Video
	for _, value := range videoList {
		if value.Id == id {
			video = value
			break
		}
	}

	b, err := json.Marshal(video)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err = io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {
	fileReader, header, err := r.FormFile("file[]")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType != videoContentType {
		http.Error(w, "Unexpected content type", http.StatusBadRequest)
		return
	}

	fileName := header.Filename

	file, err := createFile(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close() // случайно заигнорил тип ошибки unhandled error

	_, err = io.Copy(file, fileReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func createFile(fileName string) (*os.File, error) {
	if err := os.Mkdir(dirPath, os.ModeDir); err != nil && !os.IsExist(err) {
		return nil, err
	}
	filePath := filepath.Join(dirPath, fileName)
	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
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
