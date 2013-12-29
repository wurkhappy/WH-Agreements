package handlers

import (
	"encoding/json"
	"fmt"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Agreements/models"
	"github.com/wurkhappy/WH-Config"
	"net/http"
)

type StatusData struct {
	Action          string `json:"action"`
	Message         string `json:"message"`
	UserID          string `json:"userID"`
	CreditSourceURI string `json:"creditSourceURI"`
	DebitSourceURI  string `json:"debitSourceURI"`
	PaymentType     string `json:"paymentType"`
}

type Comment struct {
	UserID             string `json:"userID"`
	AgreementID        string `json:"agreementID"`
	AgreementVersionID string `json:"agreementVersionID"`
	Text               string `json:"text"`
	MilestoneID        string `json:"milestoneID"`
	StatusID           string `json:"statusID"`
}

func CreateAgreementStatus(params map[string]interface{}, body []byte) ([]byte, error, int) {
	versionID := params["versionID"].(string)
	agreement, err := models.FindAgreementByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	var data *StatusData
	json.Unmarshal(body, &data)
	status := models.CreateStatus(agreement.AgreementID, agreement.VersionID, "", data.Action, agreement.Version)
	status.UserID = data.UserID

	if agreement.CurrentStatus != nil && status.Action == agreement.CurrentStatus.Action {
		return nil, fmt.Errorf("%s", "Action already taken"), http.StatusConflict
	}

	switch status.Action {
	case "submitted":
		agreement.Draft = false
		agreement.Version += 1
		err = agreement.ArchiveOtherVersions()
		if err != nil {
			return nil, fmt.Errorf("%s %s", "Error archiving: ", err.Error()), http.StatusBadRequest
		}
		go emailSubmittedAgreement(status.AgreementID, data.Message)
	case "accepted":
		go emailAcceptedAgreement(status.AgreementID, data.Message)
	case "rejected":
		go emailRejectedAgreement(status.AgreementID, data.Message)
	}

	if data.Message != "" && data.Message != " " {
		comment := &Comment{AgreementVersionID: agreement.VersionID, Text: data.Message, StatusID: status.ID, UserID: data.UserID, AgreementID: agreement.AgreementID}
		commentBytes, _ := json.Marshal(comment)

		go sendServiceRequest("POST", config.CommentsService, "/agreement/"+agreement.AgreementID+"/comments?sendEmail=false", commentBytes)
	}

	agreement.CurrentStatus = status
	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}
	err = status.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}
	s, _ := json.Marshal(status)
	return s, nil, http.StatusOK
}

func CreatePaymentStatus(params map[string]interface{}, body []byte) ([]byte, error, int) {
	versionID := params["versionID"].(string)
	agreement, err := models.FindAgreementByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	paymentID := params["paymentID"].(string)
	payment := agreement.WorkItems.GetWorkItem(paymentID)

	var data *StatusData
	json.Unmarshal(body, &data)
	status := models.CreateStatus(agreement.AgreementID, agreement.VersionID, paymentID, data.Action, agreement.Version)
	status.UserID = data.UserID

	if payment.CurrentStatus != nil && status.Action == payment.CurrentStatus.Action {
		return nil, fmt.Errorf("%s", "Action already taken"), http.StatusConflict
	}

	switch status.Action {
	case "submitted":
		if payment.CurrentStatus != nil && payment.CurrentStatus.Action == "accepted" {
			return nil, fmt.Errorf("%s", "Action already accepted"), http.StatusConflict
		}
		go createNewTransaction(versionID, paymentID, data.CreditSourceURI)
		go emailSubmittedPayment(versionID, paymentID, data.Message)
	case "accepted":
		go sendPayment(status, data.DebitSourceURI, data.PaymentType)
		go emailSentPayment(versionID, paymentID, data.Message)
		go emailAcceptedPayment(versionID, paymentID, data.Message)
	case "rejected":
		go emailRejectedPayment(versionID, paymentID, data.Message)
	}

	payment.CurrentStatus = status
	if agreement.CurrentStatus != nil && !(status.Action == "submitted" && agreement.CurrentStatus.Action == "submitted" && agreement.CurrentStatus.ParentID == "") {
		//we're checking here if it's a deposit request. If it's not then we can update the agreement status
		//If it is a deposit request we want the agreement to keep it's submitted status because it needs to be accepted
		//before the payment is accepted
		agreement.CurrentStatus = status
	}

	//check if there's any message attached
	if data.Message != "" && data.Message != " " {
		comment := &Comment{AgreementVersionID: agreement.VersionID, Text: data.Message, StatusID: status.ID, UserID: data.UserID, AgreementID: agreement.AgreementID, MilestoneID: paymentID}
		commentBytes, _ := json.Marshal(comment)

		go sendServiceRequest("POST", config.CommentsService, "/agreement/"+agreement.AgreementID+"/comments?sendEmail=false", commentBytes)
	}

	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}
	err = status.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}
	s, _ := json.Marshal(status)
	return s, nil, http.StatusOK
}

func createNewTransaction(versionID, paymentID, creditURI string) {
	agreement, _ := models.FindAgreementByVersionID(versionID)

	var amount int
	for _, payment := range agreement.WorkItems {
		if payment.ID == paymentID {
			amount = payment.Amount
		}
	}

	m := map[string]interface{}{
		"creditSourceURI": creditURI,
		"clientID":        agreement.ClientID,
		"freelancerID":    agreement.FreelancerID,
		"agreementID":     agreement.AgreementID,
		"paymentID":       paymentID,
		"amount":          amount,
	}
	bodyJSON, _ := json.Marshal(m)
	message := map[string]interface{}{
		"Method": "POST",
		"Body":   bodyJSON,
	}

	body, _ := json.Marshal(message)
	publisher, _ := rbtmq.NewPublisher(connection, config.TransactionsExchange, "direct", config.TransactionsQueue, "/transactions")
	publisher.Publish(body, true)
}

func sendPayment(status *models.Status, debitURI string, paymentType string) {

	m := map[string]interface{}{
		"debitSourceURI": debitURI,
		"paymentType":    paymentType,
	}
	bodyJSON, _ := json.Marshal(m)
	message := map[string]interface{}{
		"Method": "PUT",
		"Body":   bodyJSON,
	}

	body, _ := json.Marshal(message)
	publisher, _ := rbtmq.NewPublisher(connection, config.TransactionsExchange, "direct", config.TransactionsQueue, "/payment/"+status.ParentID+"/transaction")
	publisher.Publish(body, true)
}
