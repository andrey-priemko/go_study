package provider

import (
	"errors"
	"go_study/videoserver/model"
)

type MockProvider struct {
	Videos []model.Video
	VideosForList []model.VideoListItem
}

func (p *MockProvider) GetVideo(id string) (*model.Video, error) {
	for _, item := range p.Videos {
		if item.Id == id {
			return &item, nil
		}
	}
	return nil, errors.New("not found")
}

func (p *MockProvider) GetVideoList() ([]model.VideoListItem, error) {
	return p.VideosForList, nil
}

func (p *MockProvider) UploadVideo(id string, fileName string, url string) error {
	return errors.New("not implemented")
}

func (p *MockProvider) Init() {
	p.Videos = append(p.Videos, model.Video{
		"id1",
		"name1",
		1,
		"thumb11",
		"url1",
	})
	p.Videos = append(p.Videos, model.Video{
		"id2",
		"name2",
		2,
		"thumb2",
		"url2",
	})

	p.VideosForList = append(p.VideosForList, model.VideoListItem{
		"id1",
		"name1",
		1,
		"thumb11",
	})
	p.VideosForList = append(p.VideosForList, model.VideoListItem{
		"id2",
		"name2",
		2,
		"thumb2",
	})
}