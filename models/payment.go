package models

import (
	"github.com/nu7hatch/gouuid"
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

type Payments []*Payment

func (p Payments) AddIDs() {
	for _, payment := range p {
		if payment.ID == "" {
			id, _ := uuid.NewV4()
			payment.ID = id.String()
		}
	}
}

func (p Payments) GetPayment(id string) *Payment {
	for _, payment := range p {
		if payment.ID == id {
			return payment
		}
	}
	return nil
}

func (p Payments) AreCompleted() bool {
	numberOfPayments := len(p)
	var numberOfPaidPayments int
	for _, payment := range p {
		if payment.CurrentStatus != nil && payment.CurrentStatus.Action == "accepted" {
			numberOfPaidPayments += 1
		}
	}

	return numberOfPayments == numberOfPaidPayments
}
