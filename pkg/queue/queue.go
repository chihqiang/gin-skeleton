package queue

import (
	"context"
	"time"
)

type ITask interface {
	Execute(val context.Context, queue IQueue)
}

type IQueue interface {
	Register(task ITask) error
	Push(task ITask, delay time.Duration) error
	Start(val context.Context, interval time.Duration)
	Stop()
}
