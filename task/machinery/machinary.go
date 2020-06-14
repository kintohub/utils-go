package machinery

import (
	"github.com/RichardKnop/machinery/v1"
	machineryConfig "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/kintohub/utils-go/config"
	"github.com/kintohub/utils-go/klog"
	"github.com/kintohub/utils-go/task"
)

type MachineryTaskClient struct {
	server *machinery.Server
}

func NewMachineryTaskClient() task.TaskClientInterface {
	var cnf = &machineryConfig.Config{
		Broker:          config.GetStringOrDie("MACHINERY_REDIS_HOST"),
		DefaultQueue:    "machinery_tasks",
		ResultBackend:   config.GetString("MACHINERY_MONGODB_HOST", ""),               // optional status
		ResultsExpireIn: config.GetInt("MACHINERY_MONGODB_EXPIRE_TIME_SECONDS", 3600), // cleanup status
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		klog.PanicfWithError(err, "could not start machinery server")
	}

	if config.GetBool("MACHINERY_WORKERS_ENABLED", false) {
		worker := server.NewWorker(
			config.GetStringOrDie("MACHINERY_WORKER_NAME"),
			config.GetInt("MACHINERY_WORKER_CONCURRENCY_LIMIT", 0), // 0 == no limit
		)
		go func() {
			// Blocking func
			err := worker.Launch()
			if err != nil {
				klog.PanicfWithError(err, "could not start worker(s)")
			}
		}()
	}

	return &MachineryTaskClient{
		server: server,
	}
}

func (m *MachineryTaskClient) SubmitTask(task *task.Task) error {
	_, err := m.server.SendTask(&tasks.Signature{
		Name: task.Name,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: task.Data,
			},
		},
		RetryCount:   config.GetInt("MACHINERY_RETRY_ATTEMPTS", 100),
		RetryTimeout: config.GetInt("MACHINERY_RETRY_TIMEOUT", 0), // 0 == fib sequence
	})

	return err
}

func (m *MachineryTaskClient) RegisterTaskHandler(taskName string, taskHandler task.TaskHandler) error {
	return m.server.RegisterTask(taskName, taskHandler)
}

func (m *MachineryTaskClient) RegisterChainTaskHandler(taskName string, chainTaskHandler task.ChainTaskHandler) error {
	return m.server.RegisterTask(taskName, func(json string) error {
		nextTask, err := chainTaskHandler(json)

		if err != nil {
			return err
		}

		return m.SubmitTask(nextTask)
	})
}
