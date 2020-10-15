package handler

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/starrybarry/schedule/pkg/scheduler"
)

type extractor struct{}

func newExtractor() extractor {
	return extractor{}
}

func (extractor) ExtractTaskID(r *http.Request) (int, error) {
	return 0, nil
}

func (extractor) ExtractTask(r *http.Request) (scheduler.Task, error) {
	task := scheduler.Task{}

	defer r.Body.Close()

	err := jsoniter.ConfigCompatibleWithStandardLibrary.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		return scheduler.Task{}, err
	}

	return task, nil
}
