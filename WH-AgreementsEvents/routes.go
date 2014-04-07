package main

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/wurkhappy/WH-Agreements/handlers"
)

var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "payment.submitted",
			Dest:    handlers.PaymentSubmitted,
		},
		urlrouter.Route{
			PathExp: "payment.accepted",
			Dest:    handlers.PaymentAccepted,
		},
	},
}
