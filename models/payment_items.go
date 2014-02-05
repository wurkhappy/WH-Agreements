package models

type PaymentItem struct {
	WorkItemID string `json:"workItemID"`
	TaskID     string `json:"taskID"`
}

type PaymentItems []*PaymentItem
