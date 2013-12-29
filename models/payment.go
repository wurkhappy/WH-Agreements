package models

import (
	"github.com/nu7hatch/gouuid"
)

type Payment struct {
	ID            string    `json:"id"`
	WorkItems     WorkItems `json:"workItems"`
	CurrentStatus *Status   `json:"currentStatus"`
}

type Payments []*Payment

func NewPayment() *Payment {
	id, _ := uuid.NewV4()
	return &Payment{
		ID: id.String(),
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
