package provider

import "go_study/videoserver/model"

type DataProvider interface {
	GetVideo(id string) (*model.Video, error)
	GetVideoList() ([]model.VideoListItem, error)
	UploadVideo(id string, fileName string, url string) error
}