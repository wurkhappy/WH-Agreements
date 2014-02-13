package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"strconv"
	"time"
)

type Payment struct {
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	DateExpected  time.Time    `json:"dateExpected"`
	PaymentItems  PaymentItems `json:"paymentItems"`
	CurrentStatus *Status      `json:"currentStatus"`
	IsDeposit     bool         `json:"isDeposit"`
	AmountDue     float64      `json:"amountDue"`
	AmountPaid    float64      `json:"amountPaid"`

	// //we won't store this but we need this data to delegate to the Payment service
	// PaymentMethodID    string `json:"paymentMethodID"`
	// RecipientAccountID string `json:"recipientAccountID"`
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

func (p *Payment) UpdateAmountDue() {
	var amount float64
	for _, item := range p.PaymentItems {
		amount += item.AmountDue
	}
	p.AmountDue = amount
}

func (p *Payment) MarshalJSON() ([]byte, error) {
	if p.ID == "" {
		id, _ := uuid.NewV4()
		p.ID = id.String()
	}
	var b bytes.Buffer
	b.WriteString(`{
		"id":"` + p.ID + `",
		"title":"` + p.Title + `",
		"dateExpected":` + p.DateExpected.Format(`"`+time.RFC3339Nano+`"`) + `,
		"isDeposit":` + strconv.FormatBool(p.IsDeposit) + `,
		"amountDue":` + fmt.Sprintf("%g", p.AmountDue) + `,
		"amountPaid":` + fmt.Sprintf("%g", p.AmountPaid) + `,
		"paymentItems":`)
	paymentItems, _ := json.Marshal(p.PaymentItems)
	b.Write(paymentItems)
	b.WriteString(`,
	"currentStatus":`)
	currentStatus, _ := json.Marshal(p.CurrentStatus)
	b.Write(currentStatus)
	b.WriteString(`}`)
	return b.Bytes(), nil
}
