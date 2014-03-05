package handlers

import (
	"encoding/json"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Agreements/models"
	"github.com/wurkhappy/WH-Config"
)

func emailSubmittedAgreement(agreement *models.Agreement, message string) {
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
	sendEmail("/agreement/submitted", body)
}

func emailAcceptedAgreement(agreement *models.Agreement, message string) {
	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
		},
	}

	body, _ := json.Marshal(payload)
	sendEmail("/agreement/accepted", body)
}

func emailRejectedAgreement(agreement *models.Agreement, message string) {

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
		},
	}

	body, _ := json.Marshal(payload)
	sendEmail("/agreement/rejected", body)
}

func emailChangedAgreement(agreement *models.Agreement, message string) {

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
		},
	}

	body, _ := json.Marshal(payload)
	sendEmail("/agreement/updated", body)
}

func emailSubmittedPayment(agreement *models.Agreement, payment *models.Payment, message string) {

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
			"payment":   payment,
		},
	}

	body, _ := json.Marshal(payload)
	sendEmail("/payment/submitted", body)

}

func emailSentPayment(agreement *models.Agreement, payment *models.Payment, message string) {

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
			"payment":   payment,
		},
	}

	body, _ := json.Marshal(payload)
	sendEmail("/payment/sent", body)

}

func emailRejectedPayment(agreement *models.Agreement, payment *models.Payment, message string) {

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
			"payment":   payment,
		},
	}

	body, _ := json.Marshal(payload)
	sendEmail("/payment/rejected", body)
}

func emailAcceptedPayment(agreement *models.Agreement, payment *models.Payment, message string) {

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
			"payment":   payment,
		},
	}

	body, _ := json.Marshal(payload)
	sendEmail("/payment/accepted", body)
}

func emailArchivedAgreement(a *models.Agreement) {
	//TODO:Fill in functionality
}

func sendEmail(path string, body []byte) {
	publisher, err := rbtmq.NewPublisher(connection, config.EmailExchange, "topic", config.EmailQueue, path)
	if err != nil {
		dialRMQ()
		publisher, _ = rbtmq.NewPublisher(connection, config.EmailExchange, "topic", config.EmailQueue, path)
	}
	publisher.Publish(body, true)
}
