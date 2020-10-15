package handler

import (
	"net/http"

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
	return scheduler.Task{}, nil
}
