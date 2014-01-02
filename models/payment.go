package models

import (
	"github.com/nu7hatch/gouuid"
	"time"
)

type Payment struct {
	ID              string       `json:"id"`
	PaymentItems    PaymentItems `json:"paymentItems"`
	CurrentStatus   *Status      `json:"currentStatus"`
	IncludesDeposit bool         `json:"includesDeposit"`
	DateCreated     time.Time    `json:"dateCreated"`

	// //we won't store this but we need this data to delegate to the Payment service
	// PaymentMethodID    string `json:"paymentMethodID"`
	// RecipientAccountID string `json:"recipientAccountID"`
}

type Payments []*Payment

func NewPayment() *Payment {
	id, _ := uuid.NewV4()
	return &Payment{
		ID:          id.String(),
		DateCreated: time.Now().UTC(),
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
