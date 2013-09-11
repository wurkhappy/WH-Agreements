package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Agreements/models"
	"net/http"
	"time"
)

func CreateAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	agreement := models.NewAgreement()

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &agreement)

	agreement.SaveAgreementWithCtx(ctx)

	a, _ := json.Marshal(agreement)
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
	json.Unmarshal(buf.Bytes(), &a)

	agreement, _ := models.FindAgreementByID(a.ID, ctx)

	json.Unmarshal(buf.Bytes(), &agreement)
	agreement.LastModified = time.Now()

	agreement.SaveAgreementWithCtx(ctx)

	jsonString, _ := json.Marshal(agreement)
	w.Write(jsonString)

}

func DeleteAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	models.DeleteAgreementWithID(id, ctx)

	fmt.Fprint(w, "Deleted User")

}
