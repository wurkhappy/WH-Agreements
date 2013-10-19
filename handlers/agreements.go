package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Agreements/models"
	"net/http"
)

func CreateAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	agreement := models.NewAgreement()

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	reqBytes := buf.Bytes()
	json.Unmarshal(reqBytes, &agreement)
	agreement.AppendStatus(models.StatusCreated)

	_ = agreement.SaveAgreementWithCtx(ctx)

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

func FindAgreements(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	var displayData []byte
	req.ParseForm()

	if userIDs, ok := req.Form["userID"]; ok {
		userID := userIDs[0]

		usersAgrmnts, _ := models.FindAgreementByClientID(userID, ctx)
		freelancerAgrmnts, _ := models.FindAgreementByFreelancerID(userID, ctx)
		usersAgrmnts = append(usersAgrmnts, freelancerAgrmnts...)

		displayData, _ = json.Marshal(usersAgrmnts)
	}

	w.Write(displayData)

}

func UpdateAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	reqData := make(map[string]interface{})
	json.Unmarshal(buf.Bytes(), &reqData)

	agreement, _ := models.FindAgreementByID(reqData["id"].(string), ctx)
	json.Unmarshal(buf.Bytes(), &agreement)

	//get the client's info
	if email, ok := reqData["clientEmail"]; ok {
		clientData := getUserInfo(email.(string))
		agreement.SetClientID(clientData["id"].(string))
	}
	_ = agreement.SaveAgreementWithCtx(ctx)

	jsonString, _ := json.Marshal(agreement)
	w.Write(jsonString)

}

func DeleteAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	vars := mux.Vars(req)
	id := vars["id"]
	models.DeleteAgreementWithID(id, ctx)

	fmt.Fprint(w, "Deleted User")

}


