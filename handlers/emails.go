package handlers

import (
	// "bytes"
	"encoding/json"
	// "fmt"
	// "github.com/gorilla/mux"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Agreements/models"
	"github.com/wurkhappy/WH-Config"
	// "log"
	// "net/http"
)

func emailSubmittedAgreement(agreementID, message string) {
	agreement, _ := models.FindAgreementByVersionID(agreementID)
	if agreement.Version > 1 {
		emailChangedAgreement(agreement, message)
	} else {
		emailNewAgreement(agreement, message)
	}
}

func emailNewAgreement(agreement *models.Agreement, message string) {

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, config.EmailExchange, "direct", config.EmailQueue, "/agreement/submitted")
	publisher.Publish(body, true)
}

func emailAcceptedAgreement(agreementID, message string) {
	agreement, _ := models.FindAgreementByVersionID(agreementID)

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, config.EmailExchange, "direct", config.EmailQueue, "/agreement/accepted")
	publisher.Publish(body, true)
}

func emailRejectedAgreement(agreementID, message string) {
	agreement, _ := models.FindAgreementByVersionID(agreementID)

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, config.EmailExchange, "direct", config.EmailQueue, "/agreement/rejected")
	publisher.Publish(body, true)
}

func emailChangedAgreement(agreement *models.Agreement, message string) {

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, config.EmailExchange, "direct", config.EmailQueue, "/agreement/updated")
	publisher.Publish(body, true)
}

func emailSubmittedPayment(agreementID, paymentID, message string) {
	agreement, _ := models.FindAgreementByVersionID(agreementID)
	var payment *models.Payment
	for _, pymnt := range agreement.Payments {
		if pymnt.ID == paymentID {
			payment = pymnt
		}
	}

	if payment.Required {
		return
	}

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
			"payment":   payment,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, config.EmailExchange, "direct", config.EmailQueue, "/payment/submitted")
	publisher.Publish(body, true)
}

func emailSentPayment(agreementID, paymentID, message string) {
	agreement, _ := models.FindAgreementByVersionID(agreementID)
	var payment *models.Payment
	for _, pymnt := range agreement.Payments {
		if pymnt.ID == paymentID {
			payment = pymnt
		}
	}

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
			"payment":   payment,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, config.EmailExchange, "direct", config.EmailQueue, "/payment/sent")
	publisher.Publish(body, true)
}

func emailRejectedPayment(agreementID, paymentID, message string) {
	agreement, _ := models.FindAgreementByVersionID(agreementID)
	var payment *models.Payment
	for _, pymnt := range agreement.Payments {
		if pymnt.ID == paymentID {
			payment = pymnt
		}
	}

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
			"payment":   payment,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, config.EmailExchange, "direct", config.EmailQueue, "/payment/rejected")
	publisher.Publish(body, true)
}

func emailAcceptedPayment(agreementID, paymentID, message string) {
	agreement, _ := models.FindAgreementByVersionID(agreementID)
	var payment *models.Payment
	for _, pymnt := range agreement.Payments {
		if pymnt.ID == paymentID {
			payment = pymnt
		}
	}

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
			"payment":   payment,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, config.EmailExchange, "direct", config.EmailQueue, "/payment/accepted")
	publisher.Publish(body, true)
}

func emailArchivedAgreement(a *models.Agreement) {
	//TODO:Fill in functionality
}
