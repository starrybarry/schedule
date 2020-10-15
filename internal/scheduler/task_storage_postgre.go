package scheduler

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/starrybarry/schedule/pkg/scheduler"
	"go.uber.org/zap"
)

type TaskStorage struct {
	pool *pgxpool.Pool

	log *zap.Logger
}

func NewTaskStorage(pool *pgxpool.Pool, log *zap.Logger) *TaskStorage {
	return &TaskStorage{pool: pool, log: log}
}

func (ts *TaskStorage) AddTask(ctx context.Context, task scheduler.Task) error {
	ts.log.Info("add task", zap.Any("task", task))

	conn, err := ts.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connect, error: %w", err)
	}

	row := conn.QueryRow(ctx,
		"INSERT into `tasks`"+
			"VALUE ($1) RETURNING id", task.Name)

	var id uint64

	if err = row.Scan(&id); err == nil {
		ts.log.Info("add task complete", zap.Uint64("id", id))
		return nil
	}

	return fmt.Errorf("row scan, error: %w", err)
}

func (ts *TaskStorage) DeleteTask(ctx context.Context, task scheduler.Task) error {
	ts.log.Info("delete task", zap.Any("task", task))

	conn, err := ts.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connect, error: %w", err)
	}

	row := conn.QueryRow(ctx,
		"DELETE FROM `tasks`"+
			"WHERE id = $1", task.ID)

	var id uint64

	if err = row.Scan(&id); err == nil {
		ts.log.Info("delete task complete", zap.Uint64("id", id))
		return nil
	}

	return fmt.Errorf("row scan, error: %w", err)
}
