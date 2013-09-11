package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kr/s3"
	"github.com/wurkhappy/WH-UserService/DB"
	"github.com/wurkhappy/WH-UserService/models"
	"log"
	"net/http"
	"time"
)

func CreateAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	agreement := models.NewAgreement()

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &user)

	agreement.SaveUserWithCtx(ctx)

	u, _ := json.Marshal(user)
	w.Write(u)
}

func GetAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	agreement, _ := models.FindAgreementByID(id, ctx)

	u, _ := json.Marshal(agreement)
	w.Write(u)

}

func UpdateUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	//unmarshal json into agreement struct so we have easy access to the ID field
	//with ID we pull the agreement from the DB, update it with the json and then set it's date modified field
	a := new(Agreement)

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	json.Unmarshal(buf.Bytes(), &a)

	agreement, _ := models.FindAgreementByID(a.ID, ctx)

	json.Unmarshal(buf.Bytes(), &agreement)
	agreement.LastModified = time.Now()

	agreement.SaveUserWithCtx(ctx)

	jsonString, _ := json.Marshal(agreement)
	w.Write(jsonString)

}

func DeleteUser(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	models.DeleteAgreementWithID(id, ctx)

	fmt.Fprint(w, "Deleted User")

}
