package models

import (
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

func test_WorkItemsAreCompleted(t *testing.T) {
	agreement := NewAgreement()
	workItem1 := new(WorkItem)
	workItem1.CurrentStatus = CreateStatus(agreement.AgreementID, agreement.VersionID, "", "accepted", agreement.Version)
	agreement.WorkItems = append(agreement.WorkItems, workItem1)
	if !agreement.WorkItems.AreCompleted() {
		t.Error("incomplete workItems when workItems is completed")
	}
	workItem2 := new(WorkItem)
	agreement.WorkItems = append(agreement.WorkItems, workItem2)
	if agreement.WorkItems.AreCompleted() {
		t.Error("completed workItems when workItems is incomplete")
	}
}
