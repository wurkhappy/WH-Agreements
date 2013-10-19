package handlers

import (
	// "bytes"
	"encoding/json"
	// "fmt"
	// "github.com/gorilla/mux"
	"github.com/streadway/amqp"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Agreements/models"
	// "log"
	// "net/http"
)

var connection *amqp.Connection

func init() {
	uri := "amqp://guest:guest@localhost:5672/"
	cn, err := amqp.Dial(uri)
	if err != nil {
		panic(err)
	}
	connection = cn
}

func emailSubmittedAgreement(agreementID, message string) {
	ctx, _ := DB.NewContext()
	agreement, _ := models.FindAgreementByID(agreementID, ctx)
	if agreement.Version > 1 {
		emailChangedAgreement(agreement, message)
	} else {
		emailNewAgreement(agreement)
	}
}

func emailNewAgreement(agreement *models.Agreement) {

	message := map[string]interface{}{
		"Body": agreement,
	}

	body, _ := json.Marshal(message)
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/agreement/submitted")
	publisher.Publish(body, true)
}

func emailAcceptedAgreement(agreementID, message string) {
	ctx, _ := DB.NewContext()
	agreement, _ := models.FindAgreementByID(agreementID, ctx)

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/agreement/accepted")
	publisher.Publish(body, true)
}

func emailRejectedAgreement(agreementID, message string) {
	ctx, _ := DB.NewContext()
	agreement, _ := models.FindAgreementByID(agreementID, ctx)

	payload := map[string]interface{}{
		"Body": map[string]interface{}{
			"agreement": agreement,
			"message":   message,
		},
	}

	body, _ := json.Marshal(payload)
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/agreement/rejected")
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
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/agreement/updated")
	publisher.Publish(body, true)
}

func emailSubmittedPayment(agreementID, paymentID, message string) {
	ctx, _ := DB.NewContext()
	agreement, _ := models.FindAgreementByID(agreementID, ctx)
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
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/payment/submitted")
	publisher.Publish(body, true)
}

func emailSentPayment(agreementID, paymentID, message string) {
	ctx, _ := DB.NewContext()
	agreement, _ := models.FindAgreementByID(agreementID, ctx)
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
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/payment/sent")
	publisher.Publish(body, true)
}

func emailRejectedPayment(agreementID, paymentID, message string) {
	ctx, _ := DB.NewContext()
	agreement, _ := models.FindAgreementByID(agreementID, ctx)
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
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/payment/rejected")
	publisher.Publish(body, true)
}

func emailAcceptedPayment(agreementID, paymentID, message string) {
	ctx, _ := DB.NewContext()
	agreement, _ := models.FindAgreementByID(agreementID, ctx)
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
	publisher, _ := rbtmq.NewPublisher(connection, "email", "direct", "email", "/payment/accepted")
	publisher.Publish(body, true)
}
