package handler

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/starrybarry/schedule/pkg/scheduler"

	"go.uber.org/zap"
)

func (h *handler) AddTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	task, err := h.extractor.ExtractTask(r)
	if err != nil {
		h.failed(w, scheduler.NewResponse(false), http.StatusBadRequest)
		h.log.Error("extract task", zap.Error(err))
		return
	}

	if err := h.schedulerSV.AddTask(ctx, task); err != nil {
		h.failed(w, scheduler.NewResponse(false), http.StatusInternalServerError)
		h.log.Error("add task", zap.Error(err))
		return
	}

	b, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(scheduler.NewResponse(true))
	if err != nil {
		h.failed(w, scheduler.NewResponse(false), http.StatusInternalServerError)
		h.log.Error("marshal response", zap.Error(err))
		return
	}

	h.success(w, b)

	h.log.Info("add task complete!")
}
