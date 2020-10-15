package handler

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/starrybarry/schedule/pkg/scheduler"
	"go.uber.org/zap"
)

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id, err := h.extractor.ExtractTaskID(r)
	if err != nil {
		h.failed(w, scheduler.NewResponse(false), http.StatusBadRequest)
		h.log.Error("extract id", zap.Error(err))
		return
	}

	if err := h.schedulerSV.DeleteTask(ctx, scheduler.Task{ID: id}); err != nil {
		h.failed(w, scheduler.NewResponse(false), http.StatusInternalServerError)
		h.log.Error("delete task", zap.Error(err))
		return
	}

	b, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(scheduler.NewResponse(true))
	if err != nil {
		h.failed(w, scheduler.NewResponse(false), http.StatusInternalServerError)
		h.log.Error("marshal response", zap.Error(err))
		return
	}

	h.success(w, b)

	h.log.Info("delete task complete!")
}
