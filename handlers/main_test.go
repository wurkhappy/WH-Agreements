package handlers

import (
	"encoding/json"
	"github.com/wurkhappy/WH-Agreements/DB"
	// "github.com/wurkhappy/WH-Agreements/models"
	"github.com/wurkhappy/WH-Config"
	"math/rand"
	// "strconv"
	// "log"
	"github.com/nu7hatch/gouuid"
	"testing"
	"time"
)

func init() {
	config.Test()
	Setup()
	DB.Name = "testdb"
	DB.Setup()
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
	test_FindAgreements(t)

	// DB.DB.Exec("DELETE from agreement")
	defer DB.DB.Close()
}

func test_CreateAgreement(t *testing.T) {
	var err error
	var statusCode int
	var resp []byte
	params := map[string]interface{}{}

	agrmntBytes := createMockAgreement("", "")
	resp, err, statusCode = CreateAgreement(params, agrmntBytes)
	if err != nil {
		t.Errorf("error creating user %s", err.Error())
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
	var respData map[string]interface{}
	json.Unmarshal(resp, &respData)
	var bodyData map[string]interface{}
	json.Unmarshal(agrmntBytes, &bodyData)
	if respData["title"].(string) != bodyData["title"].(string) || respData["proposedServices"].(string) != bodyData["proposedServices"].(string) {
		t.Error("wrong info returned")
	}
}

func test_GetAgreement(t *testing.T) {
	var err error
	var statusCode int
	var resp []byte
	params := map[string]interface{}{}

	//
	params["id"] = "invalidid"
	_, err, statusCode = GetAgreement(params, []byte(""))
	if err == nil {
		t.Error("invalid id did not return error")
	}
	if statusCode < 400 {
		t.Error("wrong status code returned")
	}

	resp, err, statusCode = CreateAgreement(params, createMockAgreement("", ""))
	var respData map[string]interface{}
	json.Unmarshal(resp, &respData)
	params["id"] = respData["versionID"].(string)
	resp, err, statusCode = GetAgreement(params, []byte(""))
	if err != nil {
		t.Errorf("error getting user %s", err.Error())
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
}

func test_FindAgreements(t *testing.T) {
	var err error
	var statusCode int
	var resp []byte
	params := map[string]interface{}{}

	//
	params["id"] = "invalidid"
	resp, _, statusCode = FindUserAgreements(params, []byte(""))
	var respArray []map[string]interface{}
	err = json.Unmarshal(resp, &respArray)
	if err != nil {
		t.Errorf("error parsing json: %s", err.Error())
	}
	if len(respArray) != 0 {
		t.Error("invalid id returned incorrect number of agreements")
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}

	clientID, _ := uuid.NewV4()
	params = map[string]interface{}{}
	CreateAgreement(params, createMockAgreement("", clientID.String()))
	CreateAgreement(params, createMockAgreement("", clientID.String()))
	CreateAgreement(params, createMockAgreement(clientID.String(), ""))
	params["id"] = clientID.String()
	resp, _, statusCode = FindUserAgreements(params, []byte(""))
	respArray = []map[string]interface{}{}
	err = json.Unmarshal(resp, &respArray)
	if err != nil {
		t.Errorf("error parsing json: %s", err.Error())
	}
	if len(respArray) != 2 {
		//mock agreement returns draft set to true
		//this handler only returns live agreements for clients
		//that's why we only test for 2 instead of 3
		t.Errorf("invalid id returned incorrect number of agreements %s", respArray)
	}
	if statusCode >= 400 {
		t.Error("wrong status code returned")
	}
}
