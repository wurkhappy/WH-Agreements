package models

import (
	"github.com/nu7hatch/gouuid"
	"time"
)

//a work item is anything that has a cost attached to it
//could be a milestone, could be a specific task that has an estimate
type WorkItem struct {
	//these are attributes that the service manages the data for
	//eventually this stuff could be broken out into its own service but overkill right now
	ID             string    `json:"id"`
	Tasks          Tasks     `json:"scopeItems"`
	Title          string    `json:"title"`
	DateExpected   time.Time `json:"dateExpected"`
	Description    string    `json:"description"`
	IsPaid         bool      `json:"isPaid"`
	Completed      bool      `json:"completed"`
	HoursCompleted int       `json:"hoursCompleted"`
}

type WorkItems []*WorkItem

func (w WorkItems) AddIDs() {
	for _, workItem := range w {
		if workItem.ID == "" {
			id, _ := uuid.NewV4()
			workItem.ID = id.String()
		}
	}
}

func (w WorkItems) GetWorkItem(id string) *WorkItem {
	for _, workItem := range w {
		if workItem.ID == id {
			return workItem
		}
	}
	return nil
}

func (w WorkItems) AreCompleted() bool {

	return true
}

func (workItems WorkItems) UpdatePaidItems(payItems PaymentItems) {
	// cachedWorkItems := make(map[string]*WorkItem)
	// for _, p := range payItems {
	// 	var w *WorkItem
	// 	var ok bool
	// 	if w, ok = cachedWorkItems[p.WorkItemID]; !ok {
	// 		w = workItems.GetWorkItem(p.WorkItemID)
	// 		cachedWorkItems[p.WorkItemID] = w
	// 	}
	// 	if p.TaskID != "" {
	// 		task := w.Tasks.GetByID(p.TaskID)
	// 		task.Completed = true
	// 		task.IsPaid = true
	// 		task.HoursCompleted = p.HoursCompleted
	// 	} else {
	// 		w.IsPaid = true
	// 		w.Completed = true
	// 		w.Tasks.SetPaid()
	// 		w.HoursCompleted = p.HoursCompleted
	// 	}
	// }
	// for _, w := range workItems {
	// 	if w.TaskArePaid() {
	// 		w.IsPaid = true
	// 		w.Completed = true
	// 	}
	// }
}

func (w *WorkItem) TaskArePaid() bool {
	return w.Tasks.ArePaid()
}
