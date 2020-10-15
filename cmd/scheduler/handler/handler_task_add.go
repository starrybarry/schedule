package handler

import (
	"context"
	"net/http"

	"github.com/starrybarry/schedule/pkg/scheduler"

	"go.uber.org/zap"
)

func (h *handler) AddTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	task, err := h.extractor.ExtractTask(r)
	if err != nil {
		_ = scheduler.NewResponse(false)
		h.log.Error("extract task", zap.Error(err))
	}

	if err := h.schedulerSV.AddTask(ctx, task); err != nil {
		_ = scheduler.NewResponse(false)
		h.log.Error("add task", zap.Error(err))
	}

	_ = scheduler.NewResponse(true)

	h.log.Info("add task complete!")
}
