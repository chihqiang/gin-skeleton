package routes

import (
	"wangzhiqiang/skeleton/app/tasks"
	"wangzhiqiang/skeleton/pkg/queue"
)

func init() {
	queue.Register(&tasks.EmailTask{})
}
