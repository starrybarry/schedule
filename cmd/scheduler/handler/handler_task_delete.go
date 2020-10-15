package handler

import (
	"context"
	"net/http"

	"github.com/starrybarry/schedule/pkg/scheduler"
	"go.uber.org/zap"
)

func (h *handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id, err := h.extractor.ExtractTaskID(r)
	if err != nil {
		_ = scheduler.NewResponse(false)
		h.log.Error("extract id", zap.Error(err))
	}

	if err := h.schedulerSV.DeleteTask(ctx, scheduler.Task{ID: id}); err != nil {
		_ = scheduler.NewResponse(false)
		h.log.Error("delete task", zap.Error(err))
	}

	_ = scheduler.NewResponse(true)

	h.log.Info("delete task complete!")
}
