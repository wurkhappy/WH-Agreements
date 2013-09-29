package models

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Payment struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	Amount        float64       `json:"amount"`
	ScopeItems    []*ScopeItem  `json:"scopeItems"`
	Title         string        `json:"title"`
	StatusHistory statusHistory `json:"statusHistory"`
	DateExpected  time.Time     `json:"dateExpected"`
}

type ScopeItem struct {
	Text string `json:"text"`
}
