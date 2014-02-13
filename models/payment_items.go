package models

type PaymentItem struct {
	WorkItemID string  `json:"workItemID"`
	TaskID     string  `json:"taskID"`
	Hours      float64 `json:"hours"`
	AmountDue  float64 `json:"amountDue"`
	Rate       float64 `json:"rate"`
	Title      string  `json:"title"`
}

type PaymentItems []*PaymentItem
