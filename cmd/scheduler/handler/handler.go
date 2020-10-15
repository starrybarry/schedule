package handler

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"

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

func (h *handler) success(w http.ResponseWriter, body []byte) {
	code := http.StatusOK

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	_, err := w.Write(body)
	if err != nil {
		h.log.Error("failed writing response", zap.Error(err), zap.ByteString("body", body))
		return
	}

	h.log.Debug("send body complete: ", zap.ByteString("body", body))
}

func (h *handler) failed(w http.ResponseWriter, v interface{}, code int) {
	body, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
	if err != nil {
		h.log.Error("failed to marshal response")

		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	_, err = w.Write(body)
	if err != nil {
		h.log.Error("failed writing response", zap.Error(err), zap.ByteString("body", body))
	}

	h.log.Debug("send error complete: ", zap.ByteString("body", body))
}

func NewHttpHandler(schedulerSv schedulerSV, log *zap.Logger) http.Handler {
	newRouter := mux.NewRouter()

	newRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("routed not found", zap.String("url", r.URL.String()))
		http.NotFound(w, r)
	})

	handler := newHandler(schedulerSv, log)

	newRouter.HandleFunc("/tasks", handler.AddTask).
		Name("add_task").Methods(http.MethodPost)

	newRouter.HandleFunc("/tasks/{id}", handler.DeleteTask).
		Name("delete_task").Methods(http.MethodDelete)

	return newRouter
}
