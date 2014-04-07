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
	fmt.Println(action)

	agreement.LastAction = action

	if action.Name == models.ActionCompleted {
		agreement.Archived = true
	} else if action.Name == models.ActionSubmitted {
		agreement.Version += 1
	}

	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("Error saving agreement", err.Error()), http.StatusBadRequest
	}

	a, _ := json.Marshal(action)
	return a, nil, http.StatusOK
}
