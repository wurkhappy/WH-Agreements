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
	publisher, _ := rbtmq.NewPublisher(connection, config.EmailExchange, "direct", config.EmailQueue, "/agreement/submitted")
	publisher.Publish(body, true)
}

func emailAcceptedAgreement(agreement *models.Agreement, message string) {
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

func emailRejectedAgreement(agreement *models.Agreement, message string) {

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

func emailSubmittedPayment(agreement *models.Agreement, payment *models.Payment, message string) {

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

func emailSentPayment(agreement *models.Agreement, payment *models.Payment, message string) {

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

func emailRejectedPayment(agreement *models.Agreement, payment *models.Payment, message string) {

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

func emailAcceptedPayment(agreement *models.Agreement, payment *models.Payment, message string) {

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
