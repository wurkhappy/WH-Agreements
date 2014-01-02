package models

type PaymentItem struct {
	WorkItemID string `json:"workItemID"`
	Amount     int    `json:"amount"`
}

type PaymentItems []*PaymentItem
