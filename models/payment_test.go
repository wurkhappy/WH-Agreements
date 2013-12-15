package models

import (
	"github.com/nu7hatch/gouuid"
	"testing"
	// "time"
)

func test_AddIDs(t *testing.T) {
	payment1 := new(Payment)
	payment2 := new(Payment)

	var payments Payments
	payments = append(payments, payment1, payment2)

	payments.AddIDs()
	if payment1.ID == "" || payment2.ID == "" {
		t.Error("IDs weren't added to payments")
	}
}

func test_GetPayment(t *testing.T) {
	payment1 := new(Payment)
	id, _ := uuid.NewV4()
	payment1.ID = id.String()

	var payments Payments
	payments = append(payments, payment1)

	payment := payments.GetPayment(id.String())
	if payment != payment1 {
		t.Error("wrong payment was returned")
	}
}

func test_PaymentsAreCompleted(t *testing.T) {
	agreement := NewAgreement()
	payment1 := new(Payment)
	payment1.CurrentStatus = CreateStatus(agreement.AgreementID, agreement.VersionID, "", "accepted", agreement.Version)
	agreement.Payments = append(agreement.Payments, payment1)
	if !agreement.Payments.AreCompleted() {
		t.Error("incomplete payments when payments is completed")
	}
	payment2 := new(Payment)
	agreement.Payments = append(agreement.Payments, payment2)
	if agreement.Payments.AreCompleted() {
		t.Error("completed payments when payments is incomplete")
	}
}
