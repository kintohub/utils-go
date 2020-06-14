package task

type TaskHandler func(json string) error
type ChainTaskHandler func(json string) (*Task, error)

type TaskClientInterface interface {
	// Submit a task to your worker(s)
	SubmitTask(task *Task) error
	// Register a task handler that only processes tasks
	RegisterTaskHandler(taskName string, taskHandler TaskHandler) error
	// Register a task handler that will return a follow up task after processing its task
	RegisterChainTaskHandler(taskName string, chainTaskHandler ChainTaskHandler) error
}
