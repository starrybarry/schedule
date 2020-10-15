package scheduler

import "time"

type Task struct {
	ID       int       `json:"id" db:"id"`
	ExecTime time.Time `json:"exec_time" db:"exec_time"`
	Name     string    `json:"name" db:"name"`
}
