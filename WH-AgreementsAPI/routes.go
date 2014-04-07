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
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"PUT":    handlers.UpdateAgreement,
				"GET":    handlers.GetAgreement,
				"DELETE": handlers.DeleteAgreement,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/v/:id/action",
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"POST": handlers.UpdateAction,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/v",
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"POST": handlers.CreateAgreement,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/:id",
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"GET": handlers.GetLatestAgreement,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/:id/owners",
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"GET": handlers.GetAgreementOwner,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/agreements",
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"GET": handlers.FindUserAgreements,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/archives",
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"GET": handlers.FindUserArchivedAgreements,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/v/:id/owners",
			Dest: map[string]func(map[string]interface{}, []byte) ([]byte, error, int){
				"GET": handlers.GetVersionOwner,
			},
		},
	},
}
