package main

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/wurkhappy/WH-Agreements/handlers"
)

//order matters so most general should go towards the bottom
var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "/agreements/v/:id",
			Dest: map[string]interface{}{
				"PUT":    handlers.UpdateAgreement,
				"GET":    handlers.GetAgreement,
				"DELETE": handlers.DeleteAgreement,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/v",
			Dest: map[string]interface{}{
				"POST": handlers.CreateAgreement,
			},
		},
		urlrouter.Route{
			PathExp: "/agreement/v/:versionID/payment/:paymentID/status",
			Dest: map[string]interface{}{
				"POST": handlers.CreatePaymentStatus,
			},
		},
		urlrouter.Route{
			PathExp: "/agreement/v/:versionID/status",
			Dest: map[string]interface{}{
				"POST": handlers.CreateAgreementStatus,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/:id/owners",
			Dest: map[string]interface{}{
				"GET": handlers.GetAgreementOwner,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/agreements",
			Dest: map[string]interface{}{
				"GET": handlers.FindUserAgreements,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/v/:id/owners",
			Dest: map[string]interface{}{
				"GET": handlers.GetVersionOwner,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/v/:id/archive",
			Dest: map[string]interface{}{
				"POST": handlers.ArchiveAgreement,
			},
		},
	},
}
