package scheduler

import "time"

type Task struct {
	ID     string
	DateAt time.Time
}
