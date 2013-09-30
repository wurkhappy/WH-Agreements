package models

import (
	"time"
)

type Payment struct {
	ID            string       `json:"id" bson:"_id"`
	Amount        float64      `json:"amount"`
	ScopeItems    []*ScopeItem `json:"scopeItems"`
	Title         string       `json:"title"`
	CurrentStatus *Status      `json:"currentStatus" bson:",omitempty"`
	DateExpected  time.Time    `json:"dateExpected"`
}

type ScopeItem struct {
	Text string `json:"text"`
}
