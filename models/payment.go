package models

import (
	// "labix.org/v2/mgo/bson"
)

type Payment struct {
	// ID         bson.ObjectId `json:"id" bson:"_id"`
	Amount     float64      `json:"amount"`
	ScopeItems []*ScopeItem `json:"scopeItems"`
	Title      string       `json:"title"`
}

type ScopeItem struct {
	Text string `json:"text"`
}
