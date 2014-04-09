package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/wurkhappy/WH-Agreements/models"
	"github.com/wurkhappy/WH-Config"
	"log"
	"net/http"
	"time"
)

type Event struct {
	Name string
	Body []byte
}

type Events []*Event

func (e Events) Publish() {
	ch := getChannel()
	defer ch.Close()
	for _, event := range e {
		event.PublishOnChannel(ch)
	}
}

func (e *Event) PublishOnChannel(ch *amqp.Channel) {
	if ch == nil {
		ch = getChannel()
		defer ch.Close()
	}

	ch.ExchangeDeclare(
		config.MainExchange, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	)

	ch.Publish(
		config.MainExchange, // exchange
		e.Name,              // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        e.Body,
		})
}

func getChannel() *amqp.Channel {
	ch, err := Connection.Channel()
	if err != nil {
		dialRMQ()
		ch, err = Connection.Channel()
		if err != nil {
			log.Print(err.Error())
		}
	}

	return ch
}

func PaymentAccepted(params map[string]interface{}, body []byte) ([]byte, error, int) {
	return updatePaymentAction(body, "accepted")
}

func PaymentSubmitted(params map[string]interface{}, body []byte) ([]byte, error, int) {
	return updatePaymentAction(body, "submitted")
}

func updatePaymentAction(body []byte, actionName string) ([]byte, error, int) {
	var message struct {
		UserID    string `json:"userID"`
		VersionID string `json:"versionID"`
		Date      time.Time
	}
	json.Unmarshal(body, &message)
	action := new(models.Action)
	action.Name = actionName
	action.Type = "payment"
	action.Date = message.Date
	action.UserID = message.UserID

	var agreement *models.Agreement
	agreement, err := models.FindAgreementByVersionID(message.VersionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding agreement"), http.StatusBadRequest
	}

	agreement.LastSubAction = action

	if action.Name == "accepted" && agreement.LastAction.Name != "accepted" {
		agreement.LastAction = models.AcceptedActionForUser(message.UserID)
	}

	err = agreement.Save()
	if err != nil {
		return nil, fmt.Errorf("Error saving agreement", err.Error()), http.StatusBadRequest
	}

	return nil, nil, http.StatusOK
}
