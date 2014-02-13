package models

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
	"time"
)

type Task struct {
	ID             string    `json:"id"`
	Completed      bool      `json:"completed"`
	IsPaid         bool      `json:"isPaid"`
	HoursCompleted int       `json:"hoursCompleted"`
	Tasks          Tasks     `json:"scopeItems"`
	Title          string    `json:"title"`
	DateExpected   time.Time `json:"dateExpected"`
	Description    string    `json:"description"`
}

//for unmarshaling purposes
type task struct {
	ID             string    `json:"id"`
	Completed      bool      `json:"completed"`
	IsPaid         bool      `json:"isPaid"`
	HoursCompleted int       `json:"hoursCompleted"`
	Tasks          Tasks     `json:"scopeItems"`
	Title          string    `json:"title"`
	DateExpected   time.Time `json:"dateExpected"`
	Description    string    `json:"description"`
}

type Tasks []*Task

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
	t.Completed = tk.Completed
	t.IsPaid = tk.IsPaid
	t.HoursCompleted = tk.HoursCompleted
	t.Tasks = tk.Tasks
	t.Title = tk.Title
	t.DateExpected = tk.DateExpected
	t.Description = tk.Description
	return nil
}

func (t Tasks) GetByID(id string) *Task {
	for _, task := range t {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func (ts Tasks) ArePaid() bool {
	for _, t := range ts {
		if !t.IsPaid {
			return false
		}
	}
	return true
}
func (ts Tasks) AreComplete() bool {
	for _, t := range ts {
		if !t.Completed {
			return false
		}
	}
	return true
}
func (ts Tasks) SetPaid() {
	for _, t := range ts {
		t.IsPaid = true
		t.Completed = true
	}
}

func (ts Tasks) AddIDs() {
	for _, task := range ts {
		if task.ID == "" {
			id, _ := uuid.NewV4()
			task.ID = id.String()
		}
	}
}

func (ts Tasks) AreCompleted() bool {
	return true
}
