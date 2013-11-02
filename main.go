package main

import (
	"fmt"
	"github.com/gorilla/mux"
	// "github.com/streadway/amqp"
	// "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Agreements/handlers"
	"labix.org/v2/mgo"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", req.URL.Path[1:])
}

func main() {
	var err error
	DB.Session, err = mgo.Dial(DB.Config["DBURL"])
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/world", hello).Methods("GET")
	r.Handle("/agreement/v/{versionID}/payment/{paymentID}/status", dbContextMixIn(handlers.CreatePaymentStatus)).Methods("POST")
	r.Handle("/agreement/v/{versionID}/status", dbContextMixIn(handlers.CreateAgreementStatus)).Methods("POST")

	r.Handle("/agreements/v", dbContextMixIn(handlers.CreateAgreement)).Methods("POST")
	r.Handle("/agreements/{id}/owners", dbContextMixIn(handlers.GetAgreementOwner)).Methods("GET")
	r.Handle("/agreements", dbContextMixIn(handlers.FindAgreements)).Methods("GET")
	r.Handle("/agreements/v/{id}/owners", dbContextMixIn(handlers.GetVersionOwner)).Methods("GET")
	r.Handle("/agreements/v/{id}", dbContextMixIn(handlers.UpdateAgreement)).Methods("PUT")
	r.Handle("/agreements/v/{id}", dbContextMixIn(handlers.DeleteAgreement)).Methods("DELETE")
	r.Handle("/agreements/v/{id}", dbContextMixIn(handlers.GetAgreement)).Methods("GET")
	r.Handle("/agreements/v/{id}/archive", dbContextMixIn(handlers.ArchiveAgreement)).Methods("POST")
	http.Handle("/", r)

	http.ListenAndServe(":4050", nil)
}

type dbContextMixIn func(http.ResponseWriter, *http.Request, *DB.Context)

func (h dbContextMixIn) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//create the context
	ctx, err := DB.NewContext()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.Close()

	h(w, req, ctx)
}
