package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Agreements/models"
	"net/http"
)

func UpdateTasks(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	versionID := params["versionID"].(string)
	workItemID := params["workItemID"].(string)

	agreement, err := models.FindAgreementByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}
	// if !agreement.Draft {
	// 	return nil, fmt.Errorf("%s", "Updating not allowed"), http.StatusBadRequest
	// }

	workItem := agreement.Tasks.GetByID(workItemID)

	var tasks []*models.Task
	json.Unmarshal(body, &tasks)
	workItem.Tasks = tasks
	workItem.Completed = workItem.Tasks.AreComplete()

	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}

	t, _ := json.Marshal(tasks)
	return t, nil, http.StatusOK

}
