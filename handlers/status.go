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
	Action          string           `json:"action"`
	Message         string           `json:"message"`
	UserID          string           `json:"userID"`
	CreditSourceURI string           `json:"creditSourceURI"`
	DebitSourceURI  string           `json:"debitSourceURI"`
	PaymentType     string           `json:"paymentType"`
	WorkItems       models.WorkItems `json:"workItems"`
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

	var paymentID string
	if idpymnt, ok := params["paymentID"]; ok {
		paymentID = idpymnt.(string)
	}
	payment := agreement.Payments.GetPayment(paymentID)
	if payment == nil {
		payment = models.NewPayment()
	}

	var data *StatusData
	json.Unmarshal(body, &data)
	payment.WorkItems = data.WorkItems

	status := models.CreateStatus(agreement.AgreementID, agreement.VersionID, payment.ID, data.Action, agreement.Version)
	status.UserID = data.UserID
	payment.CurrentStatus = status

	for _, workItem := range payment.WorkItems {
		statusWI := models.CreateStatus(agreement.AgreementID, agreement.VersionID, workItem.ID, data.Action, agreement.Version)
		statusWI.UserID = data.UserID
	}

	switch status.Action {
	case "submitted":
		go createNewTransaction(agreement, payment, data.CreditSourceURI)
		go emailSubmittedPayment(versionID, payment.ID, data.Message)
	case "accepted":
		go sendPayment(payment, data.DebitSourceURI, data.PaymentType)
		go emailSentPayment(versionID, payment.ID, data.Message)
		go emailAcceptedPayment(versionID, payment.ID, data.Message)
	case "rejected":
		go emailRejectedPayment(versionID, payment.ID, data.Message)
	}

	//check if there's any message attached
	if data.Message != "" && data.Message != " " {
		comment := &Comment{AgreementVersionID: agreement.VersionID, Text: data.Message, StatusID: status.ID, UserID: data.UserID, AgreementID: agreement.AgreementID, MilestoneID: payment.ID}
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

func createNewTransaction(agreement *models.Agreement, payment *models.Payment, creditURI string) {
	var amount int
	for _, workItem := range payment.WorkItems {
		amount += workItem.Amount
	}

	m := map[string]interface{}{
		"creditSourceURI": creditURI,
		"clientID":        agreement.ClientID,
		"freelancerID":    agreement.FreelancerID,
		"agreementID":     agreement.AgreementID,
		"paymentID":       payment.ID,
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

func sendPayment(payment *models.Payment, debitURI string, paymentType string) {

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
	publisher, _ := rbtmq.NewPublisher(connection, config.TransactionsExchange, "direct", config.TransactionsQueue, "/payment/"+payment.ID+"/transaction")
	publisher.Publish(body, true)
}
