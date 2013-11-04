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

	DB.Name = "testdb"
	DB.Setup()
	DB.CreateStatements()
}

func TestIntegrationTests(t *testing.T) {
	if !testing.Short() {

		test_SaveAgreement(t)
		test_FindLatestAgreementByID(t)
		test_FindAgreementByVersionID(t)
		test_FindAgreementByFreelancerID(t)
		test_DeleteAgreementWithVersionID(t)
		test_Archive(t)
		test_ArchiveLastAgrmntVersion(t)
		test_SaveStatus(t)

		DB.DB.Exec("DELETE from agreement")
	}
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

func test_SaveAgreement(t *testing.T) {
	agreement := NewAgreement()
	err := agreement.Save()
	if err != nil {
		t.Errorf("error saving agreement: %s", err)
	}
	if time.Now().Add(time.Second * -2).After(agreement.LastModified) {
		t.Error("last modified was not updated")
	}
}

func test_FindLatestAgreementByID(t *testing.T) {
	agreement := NewAgreement()
	agreement.Save()
	a, err := FindLatestAgreementByID(agreement.AgreementID)
	if err != nil {
		t.Errorf("error finding agreement: %s", err)
	}
	if a.AgreementID != agreement.AgreementID {
		t.Error("wrong agreement was returned")
	}

	_, err = FindLatestAgreementByID("invalid id")
	if err == nil {
		t.Error("agreement returned with invalid id input")
	}
	id, _ := uuid.NewV4()
	agreement.VersionID = id.String()
	agreement.Version = 2
	agreement.Save()
	a, err = FindLatestAgreementByID(agreement.AgreementID)
	if a.Version != agreement.Version {
		t.Errorf("latest version wasn't returned %s", err)
	}

}

func test_FindAgreementByVersionID(t *testing.T) {
	agreement := NewAgreement()
	agreement.Save()
	a, err := FindAgreementByVersionID(agreement.VersionID)
	if err != nil {
		t.Errorf("error finding agreement: %s", err)
	}
	if a.VersionID != agreement.VersionID {
		t.Error("wrong agreement was returned")
	}

	_, err = FindAgreementByVersionID("invalid id")
	if err == nil {
		t.Error("agreement returned with invalid id input")
	}
}

func test_FindLiveAgreementsByClientID(t *testing.T) {
	agreement1 := NewAgreement()
	id, _ := uuid.NewV4()
	agreement1.ClientID = id.String()
	agreement1.Draft = false
	agreement1.Save()
	agreement2 := NewAgreement()
	agreement2.ClientID = id.String()
	agreement2.Draft = false
	agreement2.Save()

	agreements, err := FindLiveAgreementsByClientID(agreement1.ClientID)
	if err != nil {
		t.Errorf("error finding agreements: %s", err)
	}
	if len(agreements) != 2 {
		t.Error("wrong number of agreements returned")
	}
	if agreements[0].ClientID != agreement1.ClientID {
		t.Error("wrong agreement was returned")
	}

	agreements2, _ := FindLiveAgreementsByClientID("invalid id")
	if len(agreements2) != 0 {
		t.Error("agreements returned with invalid id input")
	}
}

func test_FindAgreementByFreelancerID(t *testing.T) {
	agreement1 := NewAgreement()
	id, _ := uuid.NewV4()
	agreement1.FreelancerID = id.String()
	agreement1.Draft = false
	agreement1.Save()
	agreement2 := NewAgreement()
	agreement2.FreelancerID = id.String()
	agreement2.Draft = false
	agreement2.Save()

	agreements, err := FindAgreementByFreelancerID(agreement1.FreelancerID)
	if err != nil {
		t.Errorf("error finding agreements: %s", err)
	}
	if len(agreements) != 2 {
		t.Error("wrong number of agreements returned")
	}
	if agreements[0].ClientID != agreement1.ClientID {
		t.Error("wrong agreement was returned")
	}

	agreements2, _ := FindAgreementByFreelancerID("invalid id")
	if len(agreements2) != 0 {
		t.Error("agreements returned with invalid id input")
	}
}

func test_DeleteAgreementWithVersionID(t *testing.T) {
	agreement := NewAgreement()
	agreement.Save()
	err := DeleteAgreementWithVersionID(agreement.VersionID)

	if err != nil {
		t.Errorf("%s--- error is:%v", "test_DeleteAgreementWithVersionID", err)
	}

	a, err := FindAgreementByVersionID(agreement.VersionID)
	if a != nil {
		t.Errorf("%s--- user was found", "test_DeleteAgreementWithVersionID")
	}

	err = DeleteAgreementWithVersionID("invalid-id")
	if err == nil {
		t.Errorf("%s--- DB deleted with invalid id", "test_DeleteAgreementWithVersionID")
	}
}

func test_Archive(t *testing.T) {
	agreement := NewAgreement()
	agreement.Archive()
	a, _ := FindAgreementByVersionID(agreement.VersionID)
	if !a.Archived {
		t.Errorf("%s--- agreement wasn't archived", "test_Archive")
	}
}

func test_ArchiveLastAgrmntVersion(t *testing.T) {
	agreement := NewAgreement()
	agreement.Save()
	vOneID := agreement.VersionID

	agreement.Version = 2
	id, _ := uuid.NewV4()
	agreement.VersionID = id.String()
	agreement.Save()

	ArchiveLastAgrmntVersion(agreement.AgreementID)
	a, _ := FindAgreementByVersionID(vOneID)
	if !a.Archived {
		t.Errorf("%s--- agreement wasn't archived", "test_ArchiveLastAgrmntVersion")
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
