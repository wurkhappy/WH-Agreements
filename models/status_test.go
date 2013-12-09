package models

import (
	"github.com/nu7hatch/gouuid"
	"testing"
	"time"
)

func test_CreateStatus(t *testing.T) {
	agreementID, _ := uuid.NewV4()
	versionID, _ := uuid.NewV4()
	paymentID, _ := uuid.NewV4()
	action := "submitted"
	versionNumber := 1
	status := CreateStatus(agreementID.String(), versionID.String(), paymentID.String(), action, versionNumber)
	if status.AgreementID != agreementID.String() {
		t.Error("wrong agreement ID returned")
	}
	if status.AgreementVersionID != versionID.String() {
		t.Error("wrong version ID returned")
	}
	if status.PaymentID != paymentID.String() {
		t.Error("wrong payment ID returned")
	}
	if status.Action != action {
		t.Error("wrong action returned")
	}
	if status.AgreementVersion != versionNumber {
		t.Error("wrong version number returned")
	}
	if status.ID == "" {
		t.Error("id not set")
	}
	if  time.Now().Unix() - status.Date.Unix() > 1{
		t.Error("wrong date set")
	}
}

func test_SaveStatus(t *testing.T) {
	agreement := NewAgreement()
	status := CreateStatus(agreement.AgreementID, agreement.VersionID, "", "submitted", agreement.Version)
	err := status.Save()
	if err != nil {
		t.Errorf("error saving status: %s", err)
	}
}

func test_GetStatusHistory(t *testing.T) {
	agreement := NewAgreement()
	status1 := CreateStatus(agreement.AgreementID, agreement.VersionID, "", "submitted", agreement.Version)
	status1.Save()
	status2 := CreateStatus(agreement.AgreementID, agreement.VersionID, "", "submitted", agreement.Version)
	status2.Save()

	statusHistory, err := GetStatusHistory(agreement.AgreementID)
	if err != nil {
		t.Errorf("error finding statusHistory: %s", err)
	}
	if len(statusHistory) != 2 {
		t.Error("wrong number of statuses returned")
	}
	if statusHistory[0].AgreementID != agreement.AgreementID {
		t.Error("wrong status was returned")
	}

	agreements2, _ := GetStatusHistory("invalid id")
	if len(agreements2) != 0 {
		t.Error("statuses returned with invalid id input")
	}

}
