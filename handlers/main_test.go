package handlers

import (
	"encoding/json"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Config"
	"math/rand"
	"testing"
	"time"
)

func init() {
	config.Test()
	Setup()
	DB.Name = "testdb"
	DB.Setup(false)
	DB.CreateStatements()
	rand.Seed(time.Now().Unix())
}

func createMockAgreement(clientid, freelancerid string) []byte {
	ag := map[string]interface{}{
		"title":            "Test",
		"proposedServices": "Blah",
		"clientID":         clientid,
		"freelancerID":     freelancerid,
	}
	a, _ := json.Marshal(ag)
	return a
}

func Test(t *testing.T) {
	test_CreateAgreement(t)
	test_GetAgreement(t)
	test_FindUserAgreements(t)
	test_FindUserArchivedAgreements(t)

	DB.DB.Exec("DELETE from agreement")
	defer DB.DB.Close()
}
