package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Agreements/DB"
	"labix.org/v2/mgo"
	"testing"
	"time"
)

var ctx *DB.Context

func init() {
	var err error
	DB.Session, err = mgo.Dial(DB.Config["DBURL"])
	if err != nil {
		panic(err)
	}
	ctx = new(DB.Context)
	ctx.Database = DB.Session.Clone().DB("TestDB")
}

func Test_NewAgreement(t *testing.T) {
	agreement := NewAgreement()

	if agreement.VersionID == "" {
		t.Error("Version ID was not set")
	}

	_, err := uuid.ParseHex(agreement.VersionID)
	if err != nil {
		t.Error("agreement ID is not a UUID")
	}

	if agreement.AgreementID == "" {
		t.Error("Version ID was not set")
	}

	_, err = uuid.ParseHex(agreement.AgreementID)
	if err != nil {
		t.Error("agreement ID is not a UUID")
	}

	if agreement.AgreementID != agreement.VersionID {
		t.Error("Agreement ID is not being set to Version ID")
	}

	if agreement.Version != 1 {
		t.Error("New agreement isn't being created. Bad version number")
	}

	if !agreement.Draft {
		t.Error("New agreement isn't being set as draft")
	}
}

func Test_AddIDtoPayments(t *testing.T) {
	agreement := NewAgreement()
	payment := new(Payment)
	payment.Title = "test"
	agreement.Payments = append(agreement.Payments, payment)
	agreement.AddIDtoPayments()
	if payment.ID == "" {
		t.Error("ID was not added to payment")
	}
}

func Test_SaveAgreementWithCtx(t *testing.T) {
	agreement := NewAgreement()
	err := agreement.SaveAgreementWithCtx(ctx)
	if err != nil {
		t.Errorf("error saving agreement: %s", err)
	}
	if time.Now().Add(time.Second * -2).After(agreement.LastModified) {
		t.Error("last modified was not updated")
	}
}

func Test_FindLatestAgreementByID(t *testing.T) {
	agreement := NewAgreement()
	agreement.SaveAgreementWithCtx(ctx)
	a, err := FindLatestAgreementByID(agreement.AgreementID, ctx)
	if err != nil {
		t.Errorf("error finding agreement: %s", err)
	}
	if a.AgreementID != agreement.AgreementID {
		t.Error("wrong agreement was returned")
	}

	_, err = FindLatestAgreementByID("invalid id", ctx)
	if err == nil {
		t.Error("agreement returned with invalid id input")
	}
}

func Test_FindAgreementByVersionID(t *testing.T) {
	agreement := NewAgreement()
	agreement.SaveAgreementWithCtx(ctx)
	a, err := FindAgreementByVersionID(agreement.VersionID, ctx)
	if err != nil {
		t.Errorf("error finding agreement: %s", err)
	}
	if a.VersionID != agreement.VersionID {
		t.Error("wrong agreement was returned")
	}

	_, err = FindAgreementByVersionID("invalid id", ctx)
	if err == nil {
		t.Error("agreement returned with invalid id input")
	}
}

func Test_FindLiveAgreementsByClientID(t *testing.T) {
	agreement1 := NewAgreement()
	id, _ := uuid.NewV4()
	agreement1.ClientID = id.String()
	agreement1.Draft = false
	agreement1.SaveAgreementWithCtx(ctx)
	agreement2 := NewAgreement()
	agreement2.ClientID = id.String()
	agreement2.Draft = false
	agreement2.SaveAgreementWithCtx(ctx)

	agreements, err := FindLiveAgreementsByClientID(agreement1.ClientID, ctx)
	if err != nil {
		t.Errorf("error finding agreements: %s", err)
	}
	if len(agreements) != 2 {
		t.Error("wrong number of agreements returned")
	}
	if agreements[0].ClientID != agreement1.ClientID {
		t.Error("wrong agreement was returned")
	}

	agreements2, _ := FindLiveAgreementsByClientID("invalid id", ctx)
	if len(agreements2) != 0 {
		t.Error("agreements returned with invalid id input")
	}
}
