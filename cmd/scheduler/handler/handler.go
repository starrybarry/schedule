package handler

import (
	"context"
	"net/http"

	"github.com/starrybarry/schedule/pkg/scheduler"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type schedulerSV interface {
	AddTask(ctx context.Context, task scheduler.Task) error
	DeleteTask(ctx context.Context, task scheduler.Task) error
}

type handler struct {
	schedulerSV schedulerSV
	extractor   extractor
	log         *zap.Logger
}

func newHandler(schedulerSv schedulerSV, log *zap.Logger) handler {
	return handler{
		extractor:   newExtractor(),
		schedulerSV: schedulerSv,
		log:         log,
	}
}

func NewHttpHandler(schedulerSv schedulerSV, log *zap.Logger) http.Handler {
	newRouter := mux.NewRouter()

	newRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("routed not found", zap.String("url", r.URL.String()))
		http.NotFound(w, r)
	})

	handler := newHandler(schedulerSv, log)

	newRouter.HandleFunc("/task", handler.AddTask).
		Name("add_task").Methods(http.MethodPost)

	newRouter.HandleFunc("/task", handler.DeleteTask).
		Name("delete_task").Methods(http.MethodDelete)

	return newRouter
}
