package provider

import "go_study/videoserver/model"

type DataProvider interface {
	GetVideo(id string) (*model.Video, error)
}