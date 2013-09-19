package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Agreements/models"
	"net/http"
	"log"
)

func CreateAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	agreement := models.NewAgreement()

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	agreement.UnmarshalJSON(buf.Bytes())
	log.Print(agreement)
	// json.Unmarshal(buf.Bytes(), &agreement)

	err := agreement.SaveAgreementWithCtx(ctx)
	log.Print(err)

	a, _ := agreement.GetJSON()
	w.Write(a)
}

func GetAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	agreement, _ := models.FindAgreementByID(id, ctx)

	u, _ := json.Marshal(agreement)
	w.Write(u)

}

func UpdateAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	//unmarshal json into agreement struct so we have easy access to the ID field
	//with ID we pull the agreement from the DB, update it with the json and then set it's date modified field
	a := new(models.Agreement)

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	a.UnmarshalJSON(buf.Bytes())

	agreement, _ := models.FindAgreementByID(a.GetID(), ctx)

	agreement.UnmarshalJSON(buf.Bytes())
	err := agreement.SaveAgreementWithCtx(ctx)
	log.Print(err)

	jsonString, _ := agreement.GetJSON()
	w.Write(jsonString)

}

func DeleteAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	models.DeleteAgreementWithID(id, ctx)

	fmt.Fprint(w, "Deleted User")

}
