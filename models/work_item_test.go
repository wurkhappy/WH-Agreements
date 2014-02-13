package models

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
	"testing"
	// "time"
)

func test_AddIDs(t *testing.T) {
	workItem1 := new(WorkItem)
	workItem2 := new(WorkItem)

	var workItems WorkItems
	workItems = append(workItems, workItem1, workItem2)

	workItems.AddIDs()
	if workItem1.ID == "" || workItem2.ID == "" {
		t.Error("IDs weren't added to workItems")
	}
}

func test_GetWorkItem(t *testing.T) {
	workItem1 := new(WorkItem)
	id, _ := uuid.NewV4()
	workItem1.ID = id.String()

	var workItems WorkItems
	workItems = append(workItems, workItem1)

	workItem := workItems.GetWorkItem(id.String())
	if workItem != workItem1 {
		t.Error("wrong workItem was returned")
	}
}

func test_UpdatePaidItems(t *testing.T) {
	workItemsJSON := `[
	{"id":"1", "scopeItems":[{"id":"1"},{"id":"2"}]},
	{"id":"2", "scopeItems":[{"id":"1"},{"id":"2"}]},
	{"id":"3", "scopeItems":[{"id":"1"},{"id":"2"}]}
	]`
	paymentItemsJSON := `[
	{"workItemID":"1", "taskID":"1"},
	{"workItemID":"1", "taskID":"2"},
	{"workItemID":"2"},
	{"workItemID":"3", "taskID":"1"}
	]`
	var workItems WorkItems
	json.Unmarshal([]byte(workItemsJSON), &workItems)
	var paymentItems PaymentItems
	json.Unmarshal([]byte(paymentItemsJSON), &paymentItems)

	workItems.UpdatePaidItems(paymentItems)
	w1 := workItems.GetWorkItem("1")
	if !w1.IsPaid {
		t.Error("work items are not being marked paid")
	}
	if !w1.Tasks[0].IsPaid {
		t.Error("tasks are not being marked paid")
	}
	w2 := workItems.GetWorkItem("2")
	if !w2.IsPaid {
		t.Error("work items are not being marked paid")
	}
	if !w2.Tasks[0].IsPaid {
		t.Error("tasks are not being marked paid")
	}
	w3 := workItems.GetWorkItem("3")
	if w3.IsPaid {
		t.Error("work items are incorrectly being marked paid")
	}
	if !w2.Tasks[0].IsPaid {
		t.Error("tasks are not being marked paid")
	}
}
