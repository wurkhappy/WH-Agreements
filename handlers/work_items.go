package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Agreements/models"
	"net/http"
)

func UpdateWorkItem(params map[string]interface{}, body []byte) ([]byte, error, int) {
	versionID := params["versionID"].(string)
	workItemID := params["workItemID"].(string)

	agreement, err := models.FindAgreementByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	workItem := agreement.WorkItems.GetWorkItem(workItemID)

	var wi *models.WorkItem
	json.Unmarshal(body, &wi)
	//only allow description to be updated directly.
	//anything else requires an agreement change
	workItem.Description = wi.Description

	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}

	w, _ := json.Marshal(workItem)
	return w, nil, http.StatusOK

}
