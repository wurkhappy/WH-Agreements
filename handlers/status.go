package handlers

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	"github.com/wurkhappy/WH-Agreements/DB"
	"github.com/wurkhappy/WH-Agreements/models"
	"net/http"
)

func CreateAgreementStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	req.ParseForm()
	status := createStatus(req)
	if status.Action == "submitted" {
		models.ArchiveLastAgrmntVersion(status.AgreementID, ctx)
	}

	status.AddAgreementStatus(ctx)
	s, _ := json.Marshal(status)
	w.Write(s)
}

func CreatePaymentStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	status := createStatus(req)

	status.AddPaymentStatus(ctx)
	s, _ := json.Marshal(status)
	w.Write(s)
}

func createStatus(req *http.Request) *models.Status {
	status := new(models.Status)
	req.ParseForm()

	vars := mux.Vars(req)
	agreementID := vars["agreementID"]
	paymentID := vars["paymentID"]

	if actions, ok := req.Form["action"]; ok {
		action := actions[0]
		switch action {
		case "created":
			status = models.StatusCreated(agreementID, paymentID)
		case "submitted":
			status = models.StatusSubmitted(agreementID, paymentID)
		case "accepted":
			status = models.StatusAccepted(agreementID, paymentID)
		case "rejected":
			status = models.StatusRejected(agreementID, paymentID)
		}
	}

	return status
}

func UpdateAgreementStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	status := new(models.Status)

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	reqBytes := buf.Bytes()
	json.Unmarshal(reqBytes, &status)

	status.UpdateAgreementStatus(ctx)

	s, _ := json.Marshal(status)
	w.Write(s)
}

func UpdatePaymentStatus(w http.ResponseWriter, req *http.Request, ctx *DB.Context) {
	status := new(models.Status)

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	reqBytes := buf.Bytes()
	json.Unmarshal(reqBytes, &status)

	status.UpdatePaymentStatus(ctx)

	s, _ := json.Marshal(status)
	w.Write(s)
}
