package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Agreements/models"
	"github.com/wurkhappy/WH-Config"
	"net/http"
)

func CreateAgreement(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
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

func GetAgreement(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
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

func FindUserAgreements(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	usrID := params["id"].(string)
	usersAgrmnts, _ := models.FindAgreementByUserID(usrID)

	displayData, _ := json.Marshal(usersAgrmnts)
	return displayData, nil, http.StatusOK
}

func FindUserArchivedAgreements(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	var usersAgrmnts []*models.Agreement
	usrID := params["id"].(string)
	clientAgrmnts, _ := models.FindArchivedByClientID(usrID)
	freelancerAgrmnts, _ := models.FindArchivedByFreelancerID(usrID)
	usersAgrmnts = append(usersAgrmnts, freelancerAgrmnts...)
	usersAgrmnts = append(usersAgrmnts, clientAgrmnts...)

	displayData, _ := json.Marshal(usersAgrmnts)
	return displayData, nil, http.StatusOK
}

func UpdateAgreement(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	id := params["id"].(string)
	var reqData struct {
		ClientEmail string `json:"clientEmail"`
	}
	json.Unmarshal(body, &reqData)
	var newAgreement *models.Agreement
	json.Unmarshal(body, &newAgreement)

	agreement, err := models.FindAgreementByVersionID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}
	if !agreement.Draft {
		return nil, fmt.Errorf("%s", "Updating not allowed"), http.StatusBadRequest
	}
	json.Unmarshal(body, &agreement)

	//we need to do this in case the order of the work items changes
	agreement.WorkItems = newAgreement.WorkItems

	if reqData.ClientEmail != "" {
		recipientResp, _ := sendServiceRequest("GET", config.UserService, "/user/search?create=true&email="+reqData.ClientEmail, nil, userID)
		var recipientsData []map[string]interface{}
		json.Unmarshal(recipientResp, &recipientsData)
		clientData := recipientsData[0]
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

func DeleteAgreement(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	id := params["id"].(string)
	err := models.DeleteAgreementWithVersionID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error deleting agreement"), http.StatusBadRequest
	}

	return nil, nil, http.StatusOK
}

func ArchiveAgreement(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
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

func GetAgreementOwner(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
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

func GetVersionOwner(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
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

func GetLatestAgreement(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	id := params["id"].(string)
	a, err := models.FindLatestAgreementByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	jsonData, _ := json.Marshal(a)
	return jsonData, nil, http.StatusOK
}
