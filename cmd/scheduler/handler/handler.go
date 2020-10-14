package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type handler struct {
	log *zap.Logger
}

func newHandler(log *zap.Logger) handler {
	return handler{log: log}
}

func NewHttpHandler(log *zap.Logger) http.Handler {
	newRouter := mux.NewRouter()

	newRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("routed not found", zap.String("url", r.URL.String()))
		http.NotFound(w, r)
	})

	handler := newHandler(log)

	newRouter.HandleFunc("/task", handler.AddTask).
		Name("add_task").Methods(http.MethodPost)

	newRouter.HandleFunc("/task", handler.DeleteTask).
		Name("delete_task").Methods(http.MethodDelete)

	return newRouter
}
