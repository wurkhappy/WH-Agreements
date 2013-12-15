package models

import (
	"github.com/nu7hatch/gouuid"
	"testing"
	"time"
)

func test_NewAgreement(t *testing.T) {
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

	if agreement.Version != 0 {
		t.Error("New agreement bad version number")
	}

	if !agreement.Draft {
		t.Error("New agreement isn't being set as draft")
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

func test_FindAgreementByUserID(t *testing.T) {
	agreement1 := NewAgreement()
	id, _ := uuid.NewV4()
	agreement1.FreelancerID = id.String()
	agreement1.Draft = false
	agreement1.Save()
	agreement2 := NewAgreement()
	agreement2.ClientID = id.String()
	agreement2.Draft = false
	agreement2.Save()

	agreements, err := FindAgreementByUserID(agreement1.FreelancerID)
	if err != nil {
		t.Errorf("error finding agreements: %s", err)
	}
	if len(agreements) != 2 {
		t.Error("wrong number of agreements returned")
	}

	agreements2, _ := FindAgreementByUserID("invalid id")
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

func test_SetDraftCreatorID(t *testing.T) {
	agreement := NewAgreement()
	id, _ := uuid.NewV4()
	agreement.FreelancerID = id.String()
	agreement.SetDraftCreatorID()
	if agreement.DraftCreatorID != id.String() {
		t.Error("freelancer wasn't set as draft creator")
	}

	agreement = NewAgreement()
	agreement.ClientID = id.String()
	agreement.SetDraftCreatorID()
	if agreement.DraftCreatorID != id.String() {
		t.Error("client wasn't set as draft creator")
	}
}

func test_SetRecipient(t *testing.T) {
	agreement := NewAgreement()
	id, _ := uuid.NewV4()
	agreement.FreelancerID = id.String()
	otherid, _ := uuid.NewV4()
	agreement.SetRecipient(otherid.String())
	if agreement.ClientID != otherid.String() {
		t.Error("correct client id was not set")
	}

	agreement = NewAgreement()
	agreement.ClientID = id.String()
	agreement.SetRecipient(otherid.String())
	if agreement.FreelancerID != otherid.String() {
		t.Error("correct freelancer id was not set")
	}
}
