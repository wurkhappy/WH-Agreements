package handlers

import (
	// "bytes"
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Agreements/models"
	"log"
	"net/http"
)

func CreateAgreementStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	agreementID := vars["agreementID"]
	reqData, _ := parseRequest(req)

	status := createStatus(agreementID, "", reqData["action"].(string))
	switch status.Action {
	case "submitted":
		models.ArchiveLastAgrmntVersion(status.AgreementID, ctx)
		go emailNewAgreement(status.AgreementID)
	case "accepted":
		var message string
		if msg, ok := reqData["message"]; ok {
			message = msg.(string)
		}
		go emailAcceptedAgreement(status.AgreementID, message)
	}

	status.AddAgreementStatus(ctx)
	s, _ := json.Marshal(status)
	w.Write(s)
}

func CreatePaymentStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	agreementID := vars["agreementID"]
	paymentID := vars["paymentID"]
	reqData, _ := parseRequest(req)
	log.Print(reqData)

	status := createStatus(agreementID, paymentID, reqData["action"].(string))
	if status.Action == "submitted" {
		context, _ := DB.NewContext()
		go createNewTransaction(status, reqData["creditSourceURI"].(string), context)
	} else if status.Action == "accepted" {
		context, _ := DB.NewContext()
		go sendPayment(status, reqData["debitSourceURI"].(string), context)
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
