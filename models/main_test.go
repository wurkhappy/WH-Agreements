package models

import (
	"github.com/wurkhappy/WH-Agreements/DB"
	"testing"
)

func init() {

	DB.Name = "testdb"
	DB.Setup(false)
	DB.CreateStatements()
}

func TestUnitTests(t *testing.T) {
	test_NewAgreement(t)
	test_CreateStatus(t)
	test_SetDraftCreatorID(t)
	test_SetRecipient(t)
	test_UpdatePaidItems(t)

	//payment tests
	test_AddIDs(t)
	test_GetWorkItem(t)
}
func TestIntegrationTests(t *testing.T) {
	if !testing.Short() {

		test_SaveAgreement(t)
		test_FindLatestAgreementByID(t)
		test_FindAgreementByVersionID(t)
		test_FindAgreementByFreelancerID(t)
		test_DeleteAgreementWithVersionID(t)
		test_Archive(t)
		test_SaveStatus(t)
		test_FindAgreementByUserID(t)

		DB.DB.Exec("DELETE from agreement")
	}
}
