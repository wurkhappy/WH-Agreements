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
	test_AddIDtoPayments(t)
	test_CreateStatus(t)
	test_IsCompleted(t)
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

		DB.DB.Exec("DELETE from agreement")
	}
}