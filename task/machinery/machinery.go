package machinery

import (
	"github.com/RichardKnop/machinery/v1"
	machineryConfig "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/kintohub/utils-go/klog"
	"github.com/kintohub/utils-go/task"
	"math"
)

type MachineryConfig struct {
	BrokerConnectionUri        string
	DefaultQueueName           string
	ResultBackendConnectionUri string
	ResultsExpireInSeconds     int
	WorkersEnabled             bool
	WorkerAlias                string
	WorkerConcurrencyLimit     int // 0 == no limit
	MaxRetryCount              int // When set to -1
	RetryTimeoutSeconds        int
}

type MachineryTaskClient struct {
	server *machinery.Server
	config *MachineryConfig
}

func NewMachineryTaskClient(config *MachineryConfig) task.TaskClientInterface {
	// set to -1 when we want to use max that retries that the system allows
	if config.MaxRetryCount == -1 {
		config.MaxRetryCount = math.MaxInt32
	}

	var cnf = &machineryConfig.Config{
		Broker:          config.BrokerConnectionUri,
		DefaultQueue:    config.DefaultQueueName,
		ResultBackend:   config.ResultBackendConnectionUri, // optional status
		ResultsExpireIn: config.ResultsExpireInSeconds,     // cleanup status
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		klog.PanicfWithError(err, "could not start machinery server")
	}

	if config.WorkersEnabled {
		worker := server.NewWorker(
			config.WorkerAlias,
			config.WorkerConcurrencyLimit,
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
		config: config,
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
		RetryCount:   m.config.MaxRetryCount,
		RetryTimeout: m.config.RetryTimeoutSeconds, // 0 == fib sequence
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
