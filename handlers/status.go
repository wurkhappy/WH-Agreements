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
	if data.Message != "" && data.Message != " " {
		comment := &Comment{AgreementVersionID: agreement.VersionID, Text: data.Message, StatusID: status.ID, UserID: data.UserID, AgreementID: agreement.AgreementID}
		commentBytes, _ := json.Marshal(comment)

		go sendServiceRequest("POST", config.CommentsService, "/agreement/"+agreement.AgreementID+"/comments", commentBytes)
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

	var data *StatusData
	json.Unmarshal(body, &data)
	status := models.CreateStatus(agreement.AgreementID, agreement.VersionID, paymentID, data.Action, agreement.Version)

	if data.Message != "" && data.Message != " " {
		comment := &Comment{AgreementVersionID: agreement.VersionID, Text: data.Message, StatusID: status.ID, UserID: data.UserID, AgreementID: agreement.AgreementID, MilestoneID: paymentID}
		commentBytes, _ := json.Marshal(comment)

		go sendServiceRequest("POST", config.CommentsService, "/agreement/"+agreement.AgreementID+"/comments", commentBytes)
	}

	switch status.Action {
	case "submitted":
		go createNewTransaction(versionID, paymentID, data.CreditSourceURI)
		go emailSubmittedPayment(versionID, paymentID, data.Message)
	case "accepted":
		go sendPayment(status, data.DebitSourceURI)
		go emailSentPayment(versionID, paymentID, data.Message)
		go emailAcceptedPayment(versionID, paymentID, data.Message)
	case "rejected":
		go emailRejectedPayment(versionID, paymentID, data.Message)
	}

	agreement.SetPaymentStatus(status)
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

func createNewTransaction(versionID, paymentID, creditURI string) {
	agreement, _ := models.FindAgreementByVersionID(versionID)

	var amount int
	for _, payment := range agreement.Payments {
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
	message := map[string]interface{}{
		"Method": "POST",
		"Body":   m,
	}

	body, _ := json.Marshal(message)
	publisher, _ := rbtmq.NewPublisher(connection, "transactions", "direct", "transactions", "/transactions")
	publisher.Publish(body, true)
}

func sendPayment(status *models.Status, debitURI string) {

	m := map[string]interface{}{
		"debitSourceURI": debitURI,
	}
	message := map[string]interface{}{
		"Method": "PUT",
		"Body":   m,
	}

	body, _ := json.Marshal(message)
	publisher, _ := rbtmq.NewPublisher(connection, "transactions", "direct", "transactions", "/payment/"+status.PaymentID+"/transaction")
	publisher.Publish(body, true)
}
