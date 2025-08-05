package queue

import (
	"context"
	"fmt"
	"reflect"
)

var (
	tasks []ITask
)

type ITask interface {
	Execute(val context.Context, queue IQueue) error
}

func GetTasks() []ITask {
	return tasks
}
func Register(task ITask) {
	tasks = append(tasks, task)
}

// GetTaskTypeName 获取任务的唯一类型标识（包路径/类型名）
// 要求 task 必须是指针类型
func GetTaskTypeName(task ITask) (string, error) {
	t := reflect.TypeOf(task)
	if t.Kind() != reflect.Ptr {
		return "", fmt.Errorf("task must be a pointer")
	}
	return t.Elem().PkgPath() + "/" + t.Elem().Name(), nil
}
