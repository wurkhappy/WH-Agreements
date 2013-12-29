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
	ID            string       `json:"id"`
	Amount        int          `json:"amount"`
	ScopeItems    []*ScopeItem `json:"scopeItems"`
	Title         string       `json:"title"`
	CurrentStatus *Status      `json:"currentStatus"`
	DateExpected  time.Time    `json:"dateExpected"`
	Required      bool         `json:"required"`
	Completed     bool         `json:"completed"`

	//we won't store this but we need this data to delegate to the Payment service
	PaymentMethodID    string `json:"paymentMethodID"`
	RecipientAccountID string `json:"recipientAccountID"`
}

type ScopeItem struct {
	Text string `json:"text"`
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
	numberOfWorkItems := len(w)
	var numberOfPaidWorkItems int
	for _, workItem := range w {
		if workItem.CurrentStatus != nil && workItem.CurrentStatus.Action == "accepted" {
			numberOfPaidWorkItems += 1
		}
	}

	return numberOfWorkItems == numberOfPaidWorkItems
}
