package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Agreements/models"
	"net/http"
	"time"
)

func UpdateAction(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	var agreement *models.Agreement
	agreement, err := models.FindAgreementByVersionID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	var action *models.Action
	json.Unmarshal(body, &action)
	action.Date = time.Now()
	action.UserID = params["userID"].(string)

	agreement.LastAction = action

	if action.Name == models.ActionCompleted {
		agreement.Archived = true
	} else if action.Name == models.ActionSubmitted || action.Name == models.ActionUpdated {
		agreement.Version += 1
		agreement.ArchiveOtherVersions()
	}

	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("Error saving agreement", err.Error()), http.StatusBadRequest
	}

	go createAndSendEvents(body, agreement)

	a, _ := json.Marshal(action)
	return a, nil, http.StatusOK
}

func createAndSendEvents(body []byte, agreement *models.Agreement) {
	var m map[string]interface{}
	json.Unmarshal(body, &m)
	data := map[string]interface{}{
		"versionID": agreement.VersionID,
		"message":   m["message"],
		"userID":    agreement.LastAction.UserID,
		"date":      agreement.LastAction.Date,
	}
	j, _ := json.Marshal(data)
	events := Events{&Event{"agreement." + agreement.LastAction.Name, j}}
	events.Publish()
}
