package handlers

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
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
	UserID      string `json:"userID"`
	AgreementID string `json:"agreementID"`
	Text        string `json:"text"`
	MilestoneID string `json:"milestoneID"`
	StatusID    string `json:"statusID"`
}

func CreateAgreementStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	agreementID := vars["agreementID"]
	reqData := parseRequest(req)
	var data *StatusData
	json.Unmarshal(reqData, &data)

	status := createStatus(agreementID, "", data.Action)
	if data.Message != "" && data.Message != " " {
		comment := &Comment{Text: data.Message, StatusID: status.ID, UserID: data.UserID, AgreementID: agreementID}
		commentBytes, _ := json.Marshal(comment)

		body := bytes.NewReader(commentBytes)
		r, _ := http.NewRequest("POST", "http://localhost:5050/agreement/"+agreementID+"/comments", body)
		go sendRequest(r)
	}

	switch status.Action {
	case "submitted":
		models.ArchiveLastAgrmntVersion(status.AgreementID, ctx)
		go emailSubmittedAgreement(status.AgreementID, data.Message)
	case "accepted":
		go emailAcceptedAgreement(status.AgreementID, data.Message)
	case "rejected":
		go emailRejectedAgreement(status.AgreementID, data.Message)
	}

	status.AddAgreementStatus(ctx)
	s, _ := json.Marshal(status)
	w.Write(s)
}

func CreatePaymentStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	agreementID := vars["agreementID"]
	paymentID := vars["paymentID"]
	reqData := parseRequest(req)
	var data *StatusData
	json.Unmarshal(reqData, &data)
	status := createStatus(agreementID, paymentID, data.Action)

	if data.Message != "" && data.Message != " " {
		comment := &Comment{Text: data.Message, StatusID: status.ID, UserID: data.UserID, AgreementID: agreementID, MilestoneID: paymentID}
		commentBytes, _ := json.Marshal(comment)

		body := bytes.NewReader(commentBytes)
		r, _ := http.NewRequest("POST", "http://localhost:5050/agreement/"+agreementID+"/comments", body)
		go sendRequest(r)
	}

	switch status.Action {
	case "submitted":
		context, _ := DB.NewContext()
		go createNewTransaction(status, data.CreditSourceURI, context)
		go emailSubmittedPayment(agreementID, paymentID, data.Message)
	case "accepted":
		context, _ := DB.NewContext()
		go sendPayment(status, data.DebitSourceURI, context)
		go emailSentPayment(agreementID, paymentID, data.Message)
		go emailAcceptedPayment(agreementID, paymentID, data.Message)
	case "rejected":
		go emailRejectedPayment(agreementID, paymentID, data.Message)
	}

	status.AddPaymentStatus(ctx)
	s, _ := json.Marshal(status)
	w.Write(s)
}

func createStatus(agreementID, paymentID, action string) *models.Status {
	var status *models.Status
	switch action {
	case "created":
		status = models.StatusCreated(agreementID, paymentID)
	case "submitted":
		status = models.StatusSubmitted(agreementID, paymentID)
	case "accepted":
		status = models.StatusAccepted(agreementID, paymentID)
	case "rejected":
		status = models.StatusRejected(agreementID, paymentID)
	}
	return status
}

func createNewTransaction(status *models.Status, creditURI string, ctx *DB.Context) {
	agreement, _ := models.FindAgreementByID(status.AgreementID, ctx)

	var amount int
	for _, payment := range agreement.Payments {
		if payment.ID == status.PaymentID {
			amount = payment.Amount
		}
	}

	m := map[string]interface{}{
		"creditSourceURI": creditURI,
		"clientID":        agreement.ClientID,
		"freelancerID":    agreement.FreelancerID,
		"agreementID":     agreement.ID,
		"paymentID":       status.PaymentID,
		"amount":          amount,
	}
	message := map[string]interface{}{
		"Method": "POST",
		"Body":   m,
	}
	uri := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(uri)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

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
	uri := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(uri)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	body, _ := json.Marshal(message)
	publisher, _ := rbtmq.NewPublisher(connection, "transactions", "direct", "transactions", "/payment/"+status.PaymentID+"/transaction")
	publisher.Publish(body, true)
}
