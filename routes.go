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
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"PUT":    handlers.UpdateAgreement,
				"GET":    handlers.GetAgreement,
				"DELETE": handlers.DeleteAgreement,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/v",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"POST": handlers.CreateAgreement,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/:id",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"GET": handlers.GetLatestAgreement,
			},
		},
		urlrouter.Route{
			PathExp: "/agreement/v/:versionID/payment/:paymentID/status",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"PUT": handlers.UpdatePaymentStatus,
			},
		},
		urlrouter.Route{
			PathExp: "/agreement/v/:versionID/payment/:paymentID",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"PUT": handlers.UpdatePayment,
			},
		},
		urlrouter.Route{
			PathExp: "/agreement/v/:versionID/payment/",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"POST": handlers.CreatePayment,
			},
		},
		urlrouter.Route{
			PathExp: "/agreement/v/:versionID/work_item/:workItemID/tasks",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"PUT": handlers.UpdateTasks,
			},
		},
		urlrouter.Route{
			PathExp: "/agreement/v/:versionID/work_item/:workItemID",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"PUT": handlers.UpdateWorkItem,
			},
		},
		urlrouter.Route{
			PathExp: "/agreement/v/:versionID/status",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"POST": handlers.CreateAgreementStatus,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/:id/owners",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"GET": handlers.GetAgreementOwner,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/agreements",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"GET": handlers.FindUserAgreements,
			},
		},
		urlrouter.Route{
			PathExp: "/user/:id/archives",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"GET": handlers.FindUserArchivedAgreements,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/v/:id/owners",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"GET": handlers.GetVersionOwner,
			},
		},
		urlrouter.Route{
			PathExp: "/agreements/v/:id/archive",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"POST": handlers.ArchiveAgreement,
			},
		},
	},
}
