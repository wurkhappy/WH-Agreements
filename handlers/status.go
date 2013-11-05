package handlers

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Agreements/models"
	// "log"
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

func CreateAgreementStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	versionID := vars["versionID"]
	agreement, _ := models.FindAgreementByVersionID(versionID)

	reqData := parseRequest(req)
	var data *StatusData
	json.Unmarshal(reqData, &data)
	status := models.CreateStatus(agreement.AgreementID, agreement.VersionID, "", data.Action, agreement.Version)
	if data.Message != "" && data.Message != " " {
		comment := &Comment{AgreementVersionID: agreement.VersionID, Text: data.Message, StatusID: status.ID, UserID: data.UserID, AgreementID: agreement.AgreementID}
		commentBytes, _ := json.Marshal(comment)

		body := bytes.NewReader(commentBytes)
		r, _ := http.NewRequest("POST", CommentsService+"/agreement/"+agreement.AgreementID+"/comments", body)
		go sendRequest(r)
	}

	switch status.Action {
	case "submitted":
		agreement.Draft = false
		models.ArchiveLastAgrmntVersion(status.AgreementID)
		go emailSubmittedAgreement(status.AgreementID, data.Message)
	case "accepted":
		go emailAcceptedAgreement(status.AgreementID, data.Message)
	case "rejected":
		go emailRejectedAgreement(status.AgreementID, data.Message)
	}

	agreement.CurrentStatus = status
	agreement.Save()
	status.Save()
	s, _ := json.Marshal(status)
	w.Write(s)
}

func CreatePaymentStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	versionID := vars["versionID"]
	agreement, _ := models.FindAgreementByVersionID(versionID)

	paymentID := vars["paymentID"]

	reqData := parseRequest(req)
	var data *StatusData
	json.Unmarshal(reqData, &data)
	status := models.CreateStatus(agreement.AgreementID, agreement.VersionID, paymentID, data.Action, agreement.Version)

	if data.Message != "" && data.Message != " " {
		comment := &Comment{AgreementVersionID: agreement.VersionID, Text: data.Message, StatusID: status.ID, UserID: data.UserID, AgreementID: agreement.AgreementID, MilestoneID: paymentID}
		commentBytes, _ := json.Marshal(comment)

		body := bytes.NewReader(commentBytes)
		r, _ := http.NewRequest("POST", CommentsService+"/agreement/"+agreement.AgreementID+"/comments", body)
		go sendRequest(r)
	}

	switch status.Action {
	case "submitted":
		context, _ := DB.NewContext()
		go createNewTransaction(versionID, paymentID, data.CreditSourceURI, context)
		go emailSubmittedPayment(versionID, paymentID, data.Message)
	case "accepted":
		context, _ := DB.NewContext()
		go sendPayment(status, data.DebitSourceURI, context)
		go emailSentPayment(versionID, paymentID, data.Message)
		go emailAcceptedPayment(versionID, paymentID, data.Message)
	case "rejected":
		go emailRejectedPayment(versionID, paymentID, data.Message)
	}

	agreement.SetPaymentStatus(status)
	agreement.CurrentStatus = status
	agreement.Save()
	status.Save()
	s, _ := json.Marshal(status)
	w.Write(s)
}

func createNewTransaction(versionID, paymentID, creditURI string, ctx *DB.Context) {
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

func sendPayment(status *models.Status, debitURI string, ctx *DB.Context) {

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
