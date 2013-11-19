package models

import (
	"time"
)

type Payment struct {
	//these are attributes that the service manages the data for
	//eventually this stuff could be broken out into its own service but overkill right now
	ID            string       `json:"id" bson:"_id"`
	Amount        int          `json:"amount"` //cents not dollars
	ScopeItems    []*ScopeItem `json:"scopeItems"`
	Title         string       `json:"title"`
	CurrentStatus *Status      `json:"currentStatus" bson:",omitempty"`
	DateExpected  time.Time    `json:"dateExpected"`
	Required      bool         `json:"required"`

	//we won't store this but we need this data to delegate to the payment service
	PaymentMethodID    string `json:"paymentMethodID" bson:"-"`
	RecipientAccountID string `json:"recipientAccountID" bson:"-"`
}

type ScopeItem struct {
	Text string `json:"text"`
}
