package models

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

type Task struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

//for unmarshaling purposes
type task struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

func (t *Task) UnmarshalJSON(bytes []byte) (err error) {
	var tk *task
	err = json.Unmarshal(bytes, &tk)
	if err != nil {
		return err
	}

	if tk.ID == "" {
		id, _ := uuid.NewV4()
		tk.ID = id.String()
	}

	t.ID = tk.ID
	t.Text = tk.Text
	t.Completed = tk.Completed
	return nil
}
