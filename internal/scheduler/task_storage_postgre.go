package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/starrybarry/schedule/pkg/scheduler"
	"go.uber.org/zap"
)

type TaskStorage struct {
	pool      *pgxpool.Pool
	tableName string
	log       *zap.Logger
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

	defer conn.Release()

	row := conn.QueryRow(ctx,
		"INSERT into tasks (name,exec_time,created_at) "+
			"VALUES ($1, $2, $3) RETURNING id", task.Name, task.ExecTime, time.Now())

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

	defer conn.Release()

	row := conn.QueryRow(ctx,
		"DELETE FROM tasks"+
			"WHERE id = $1", task.ID)

	var id uint64

	if err = row.Scan(&id); err == nil {
		ts.log.Info("delete task complete", zap.Uint64("id", id))
		return nil
	}

	return fmt.Errorf("row scan, error: %w", err)
}

func (ts *TaskStorage) GetActualTasks(ctx context.Context) ([]scheduler.Task, error) {
	ts.log.Info("get task")

	conn, err := ts.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire connect, error: %w", err)
	}

	defer conn.Release()
	conn.Begin(ctx)
	rows, err := conn.Query(ctx,
		"SELECT id,name,exec_time FROM tasks"+
			" WHERE $1 > exec_time and in_progress = false ORDER BY exec_time", time.Now())

	if err != nil {
		return nil, fmt.Errorf("exec query, error: %w", err)
	}

	var tasks []scheduler.Task

	for rows.Next() {
		task := scheduler.Task{}

		if err := rows.Scan(&task.ID, &task.Name, &task.ExecTime); err != nil {
			return nil, fmt.Errorf("rows scan, error: %w", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (ts *TaskStorage) UpdateTask(ctx context.Context, task scheduler.Task) error {
	ts.log.Info("update task", zap.Any("task", task))

	conn, err := ts.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connect, error: %w", err)
	}

	defer conn.Release()

	row := conn.QueryRow(ctx,
		"UPDATE tasks SET (in_progress, is_finished) = ($1,$2) "+
			"WHERE id = $3 RETURNING id", task.InProgress, task.IsFinished, task.ID)

	var id uint64

	if err = row.Scan(&id); err == nil {
		ts.log.Info("update task complete", zap.Uint64("id", id))
		return nil
	}

	return fmt.Errorf("row scan, error: %w", err)
}
