package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Agreements/models"
	"log"
	"net/http"
)

func CreateAgreement(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	agreement := models.NewAgreement()

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	agreement.UnmarshalJSON(buf.Bytes())
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

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	reqData := make(map[string]interface{})
	json.Unmarshal(buf.Bytes(), &reqData)

	agreement, _ := models.FindAgreementByID(reqData["id"].(string), ctx)

	agreement.UnmarshalJSON(buf.Bytes())
	
	//get the client's info
	if email, ok := reqData["clientEmail"]; ok {
		clientData := getClientInfo(email.(string))
		agreement.SetClientID(clientData["id"].(string))
	}
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

func getClientInfo(email string) map[string]interface{} {
	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://localhost:3000/user/search?create=true&email="+email, nil)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	clientBuf := new(bytes.Buffer)
	clientBuf.ReadFrom(resp.Body)
	var clientData []map[string]interface{}
	json.Unmarshal(clientBuf.Bytes(), &clientData)
	return clientData[0]
}
