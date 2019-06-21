package handlers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go_study/videoserver/model"
	"go_study/videoserver/provider"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessfulGetVideoRequest(t *testing.T) {
	var mockProvider provider.MockProvider
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
	var mockProvider provider.MockProvider
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

