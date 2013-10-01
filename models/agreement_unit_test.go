package models

import (
	"github.com/nu7hatch/gouuid"
	"testing"
)

func Test_NewAgreement(t *testing.T) {
	agreement := NewAgreement()

	if agreement.ID == "" {
		t.Error("agreement ID was not set")
	}

	_, err := uuid.ParseHex(agreement.ID)
	if err != nil {
		t.Error("agreement ID is not a UUID")
	}

	if agreement.AgreementID != agreement.ID {
		t.Error("Agreement ID is not being set to original ID")
	}

	if agreement.Version != 1 {
		t.Error("New agreement isn't being created. Bad version number")
	}
}

func Test_GetID(t *testing.T) {
	agreement := NewAgreement()

	id := agreement.GetID()
	if id != agreement.AgreementID {
		t.Error("GetID not returning agreement id")
	}
}

func Test_SetClientID(t *testing.T) {
	agreement := NewAgreement()
	clientID := "1"
	agreement.SetClientID(clientID)
	if agreement.ClientID != clientID {
		t.Error("Client ID not being set")
	}
}

func Test_addIDtoPayments(t *testing.T) {
	agreement :=  NewAgreement()
	payment := new(Payment)
	payment.Title = "test"
	agreement.Payments = append(agreement.Payments, payment)
	agreement.addIDtoPayments()
	if payment.ID == "" {
		t.Error("ID was not added to payment")
	}
}
