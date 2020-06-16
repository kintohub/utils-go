package task

import "encoding/json"

type Task struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewTask(name string, data interface{}) (*Task, error) {
	d, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return &Task{
		Name: name,
		Data: string(d),
	}, nil
}
