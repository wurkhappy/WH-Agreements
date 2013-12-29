package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Agreements/models"
	"net/http"
)

func CreateAgreement(params map[string]interface{}, body []byte) ([]byte, error, int) {
	agreement := models.NewAgreement()

	err := json.Unmarshal(body, &agreement)
	if err != nil {
		return nil, fmt.Errorf("%s", "Wrong value types"), http.StatusBadRequest
	}
	agreement.SetDraftCreatorID()
	agreement.WorkItems.AddIDs()
	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}

	a, _ := json.Marshal(agreement)
	return a, nil, http.StatusOK

}

func GetAgreement(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	var agreement *models.Agreement
	agreement, err := models.FindAgreementByVersionID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	agreement.StatusHistory, _ = models.GetStatusHistory(agreement.AgreementID)

	a, _ := json.Marshal(agreement)
	return a, nil, http.StatusOK
}

func FindUserAgreements(params map[string]interface{}, body []byte) ([]byte, error, int) {
	userID := params["id"].(string)
	usersAgrmnts, _ := models.FindAgreementByUserID(userID)

	displayData, _ := json.Marshal(usersAgrmnts)
	return displayData, nil, http.StatusOK
}

func FindUserArchivedAgreements(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var usersAgrmnts []*models.Agreement
	userID := params["id"].(string)
	clientAgrmnts, _ := models.FindArchivedByClientID(userID)
	freelancerAgrmnts, _ := models.FindArchivedByFreelancerID(userID)
	usersAgrmnts = append(usersAgrmnts, freelancerAgrmnts...)
	usersAgrmnts = append(usersAgrmnts, clientAgrmnts...)

	displayData, _ := json.Marshal(usersAgrmnts)
	return displayData, nil, http.StatusOK
}

func UpdateAgreement(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)

	var reqData map[string]interface{}
	json.Unmarshal(body, &reqData)

	agreement, err := models.FindAgreementByVersionID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}
	if !agreement.Draft {
		return nil, fmt.Errorf("%s", "Updating not allowed"), http.StatusBadRequest
	}
	json.Unmarshal(body, &agreement)

	//get the client's info
	if email, ok := reqData["clientEmail"]; ok {
		clientData := getUserInfo(email.(string))
		agreement.SetRecipient(clientData["id"].(string))
	}
	agreement.WorkItems.AddIDs()
	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}

	jsonString, _ := json.Marshal(agreement)
	return jsonString, nil, http.StatusOK

}

func DeleteAgreement(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	err := models.DeleteAgreementWithVersionID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error deleting agreement"), http.StatusBadRequest
	}

	return nil, nil, http.StatusOK
}

func ArchiveAgreement(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)

	agreement, err := models.FindAgreementByVersionID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	//if there are payments outstanding and the user is archiving then send an email to the other user
	if agreement.WorkItems.AreCompleted() {
		go emailArchivedAgreement(agreement)
	}

	agreement.SetAsCompleted()
	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}

	jsonString, _ := json.Marshal(agreement)
	return jsonString, nil, http.StatusOK
}

func GetAgreementOwner(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	a, err := models.FindLatestAgreementByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}
	data := struct {
		ClientID   string `json:"clientID"`
		Freelancer string `json:"freelancerID"`
	}{
		a.ClientID,
		a.FreelancerID,
	}

	jsonData, _ := json.Marshal(data)
	return jsonData, nil, http.StatusOK
}

func GetVersionOwner(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	a, err := models.FindAgreementByVersionID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}
	data := struct {
		ClientID   string `json:"clientID"`
		Freelancer string `json:"freelancerID"`
	}{
		a.ClientID,
		a.FreelancerID,
	}

	jsonData, _ := json.Marshal(data)
	return jsonData, nil, http.StatusOK
}
