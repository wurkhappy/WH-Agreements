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
	IPAddress       string           `json:"ipAddress"`
	CreditSourceURI string           `json:"creditSourceID"`
	DebitSourceURI  string           `json:"debitSourceID"`
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

func CreateAgreementStatus(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	versionID := params["versionID"].(string)
	agreement, err := models.FindAgreementByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	var data *StatusData
	json.Unmarshal(body, &data)
	status := models.CreateStatus(agreement.AgreementID, agreement.VersionID, "", data.Action, agreement.Version, data.IPAddress)
	status.UserID = data.UserID

	if agreement.CurrentStatus != nil && status.Action == agreement.CurrentStatus.Action && status.ParentID == agreement.CurrentStatus.ParentID {
		return nil, fmt.Errorf("%s", "Action already taken"), http.StatusConflict
	}
	agreement.CurrentStatus = status

	switch status.Action {
	case "submitted":
		agreement.Draft = false
		agreement.Version += 1
		err = agreement.ArchiveOtherVersions()
		if err != nil {
			return nil, fmt.Errorf("%s %s", "Error archiving: ", err.Error()), http.StatusBadRequest
		}
		go emailSubmittedAgreement(agreement, data.Message)
	case "accepted":
		go emailAcceptedAgreement(agreement, data.Message)
	case "rejected":
		go emailRejectedAgreement(agreement, data.Message)
	}

	go createComment(agreement, nil, data.Message, status.UserID)

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

func CreatePayment(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	versionID := params["versionID"].(string)
	agreement, err := models.FindAgreementByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	fmt.Println(string(body))
	payment := models.NewPayment()
	json.Unmarshal(body, &payment)
	agreement.Payments = append(agreement.Payments, payment)

	var data *StatusData
	json.Unmarshal(body, &data)
	data.Action = "submitted"

	status := models.CreateStatus(agreement.AgreementID, agreement.VersionID, payment.ID, data.Action, agreement.Version, data.IPAddress)
	status.UserID = data.UserID
	payment.CurrentStatus = status

	go createNewTransaction(agreement, payment, data.CreditSourceURI)

	for _, paymentItem := range payment.PaymentItems {
		workItem := agreement.WorkItems.GetWorkItem(paymentItem.WorkItemID)
		if workItem.Required {
			payment.IncludesDeposit = true
		}
	}
	agreement.CurrentStatus = status
	go emailSubmittedPayment(agreement, payment, data.Message)

	go createComment(agreement, payment, data.Message, status.UserID)

	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}
	err = status.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}
	p, _ := json.Marshal(payment)
	return p, nil, http.StatusOK
}

func UpdatePaymentStatus(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	//get the agreeement
	versionID := params["versionID"].(string)
	agreement, err := models.FindAgreementByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	//get the payment
	var paymentID string
	if idpymnt, ok := params["paymentID"]; ok {
		paymentID = idpymnt.(string)
	}
	payment := agreement.Payments.GetPayment(paymentID)
	if payment == nil {
		return nil, fmt.Errorf("%s", "Error finding payment"), http.StatusBadRequest
	}

	//create a status
	var data *StatusData
	json.Unmarshal(body, &data)
	status := models.CreateStatus(agreement.AgreementID, agreement.VersionID, payment.ID, data.Action, agreement.Version, data.IPAddress)
	status.UserID = data.UserID
	payment.CurrentStatus = status
	agreement.CurrentStatus = status

	//check what action the status is
	switch status.Action {
	case "submitted":
		go createNewTransaction(agreement, payment, data.CreditSourceURI)
		go emailSubmittedPayment(agreement, payment, data.Message)
	case "accepted":
		//if we are accepting a payment then let's update the work items to know if they are completed or not
		for _, paymentItem := range payment.PaymentItems {
			workItem := agreement.WorkItems.GetWorkItem(paymentItem.WorkItemID)
			workItem.AmountPaid = paymentItem.Amount
		}

		go sendPayment(payment, data.DebitSourceURI, data.PaymentType)
		go emailSentPayment(agreement, payment, data.Message)
		go emailAcceptedPayment(agreement, payment, data.Message)
	case "rejected":
		go emailRejectedPayment(agreement, payment, data.Message)
	}

	go createComment(agreement, payment, data.Message, status.UserID)

	//payments and work items are part of the agreement so we need to save the agreement
	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}

	//status history for the whole agreement is saved separately
	err = status.Save()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}
	s, _ := json.Marshal(status)
	return s, nil, http.StatusOK
}

func createNewTransaction(agreement *models.Agreement, payment *models.Payment, creditURI string) {
	var amount int
	for _, paymentItem := range payment.PaymentItems {
		amount += paymentItem.Amount
	}

	m := map[string]interface{}{
		"creditSourceID": creditURI,
		"clientID":       agreement.ClientID,
		"freelancerID":   agreement.FreelancerID,
		"agreementID":    agreement.AgreementID,
		"paymentID":      payment.ID,
		"amount":         amount,
	}
	bodyJSON, _ := json.Marshal(m)
	message := map[string]interface{}{
		"Method": "POST",
		"Body":   bodyJSON,
	}

	body, _ := json.Marshal(message)
	publisher, err := rbtmq.NewPublisher(connection, config.TransactionsExchange, "direct", config.TransactionsQueue, "/transactions")
	if err != nil {
		dialRMQ()
		publisher, _ = rbtmq.NewPublisher(connection, config.TransactionsExchange, "direct", config.TransactionsQueue, "/transactions")
	}
	publisher.Publish(body, true)
}

func sendPayment(payment *models.Payment, debitURI string, paymentType string) {

	m := map[string]interface{}{
		"debitSourceID": debitURI,
		"paymentType":   paymentType,
	}
	bodyJSON, _ := json.Marshal(m)
	message := map[string]interface{}{
		"Method": "PUT",
		"Body":   bodyJSON,
	}

	body, _ := json.Marshal(message)
	publisher, err := rbtmq.NewPublisher(connection, config.TransactionsExchange, "direct", config.TransactionsQueue, "/payment/"+payment.ID+"/transaction")
	if err != nil {
		dialRMQ()
		publisher, _ = rbtmq.NewPublisher(connection, config.TransactionsExchange, "direct", config.TransactionsQueue, "/payment/"+payment.ID+"/transaction")
	}
	publisher.Publish(body, true)
}

func createComment(agreement *models.Agreement, payment *models.Payment, message string, userID string) {
	//check if there's any message attached
	if message != "" && message != " " {
		comment := &Comment{AgreementVersionID: agreement.VersionID, Text: message, UserID: userID, AgreementID: agreement.AgreementID}
		commentBytes, _ := json.Marshal(comment)

		sendServiceRequest("POST", config.CommentsService, "/agreement/"+agreement.AgreementID+"/comments?sendEmail=false", commentBytes, userID)
	}
}
