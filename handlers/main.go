package handlers

import (
	"github.com/streadway/amqp"
	"github.com/wurkhappy/WH-Config"
)

var connection *amqp.Connection

func Setup() {
	var err error
	connection, err = amqp.Dial(config.EmailURI)
	if err != nil {
		panic(err)
	}
}