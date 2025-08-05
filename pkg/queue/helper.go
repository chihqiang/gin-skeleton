package queue

var (
	tasks []ITask
)

func GetTasks() []ITask {
	return tasks
}
func Register(task ITask) {
	tasks = append(tasks, task)
}
