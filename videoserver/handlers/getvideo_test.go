package handlers

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"go_study/videoserver/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockProvider struct {
	Videos []model.Video
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
	return nil, errors.New("not implemented")
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
}

func TestSuccessfulGetVideoRequest(t *testing.T) {
	var mockProvider MockProvider
	mockProvider.Init()

	router := Router(&mockProvider)
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, err := http.Get(ts.URL + "/api/v1/video/id1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Header.Get("Content-Type"), "application/json; charset=UTF-8")
	assert.Equal(t, response.StatusCode, http.StatusOK)

	jsonStr, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	var video model.Video
	err = json.Unmarshal(jsonStr, &video)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, video.Id, "id1")
	assert.Equal(t, video.Name, "name1")
}

func TestUnsuccessfulGetVideoRequest(t *testing.T) {
	var mockProvider MockProvider
	mockProvider.Init()

	router := Router(&mockProvider)
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, err := http.Get(ts.URL + "/api/v1/video/id3")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.StatusCode, http.StatusInternalServerError)
}

