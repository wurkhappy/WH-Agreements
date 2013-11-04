package DB

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	// "log"
)

var SaveAgreement *sql.Stmt
var UpsertAgreement *sql.Stmt
var FindAgreementByVersionID *sql.Stmt
var FindLiveAgreementsByClientID *sql.Stmt
var FindAgreementByFreelancerID *sql.Stmt
var DeleteAgreement *sql.Stmt
var FindLiveVersions *sql.Stmt
var UpsertStatus *sql.Stmt
var GetStatusHistory *sql.Stmt

func CreateStatements() {
	var err error
	SaveAgreement, err = DB.Prepare("INSERT INTO agreement(id, data) VALUES($1, $2)")
	if err != nil {
		panic(err)
	}

	FindLiveVersions, err = DB.Prepare("SELECT data FROM agreement WHERE data->>'agreementID' = $1 ORDER BY data->>'version' DESC")
	if err != nil {
		panic(err)
	}

	UpsertAgreement, err = DB.Prepare("SELECT upsert_agreement($1, $2)")
	if err != nil {
		panic(err)
	}

	FindAgreementByVersionID, err = DB.Prepare("SELECT data FROM agreement WHERE id = $1")
	if err != nil {
		panic(err)
	}

	FindLiveAgreementsByClientID, err = DB.Prepare("SELECT data FROM agreement WHERE data->>'clientID' = $1 AND data->>'archived' = 'false' AND data->>'draft' = 'false'")
	if err != nil {
		panic(err)
	}

	FindAgreementByFreelancerID, err = DB.Prepare("SELECT data FROM agreement WHERE data->>'freelancerID' = $1 AND data->>'archived' = 'false'")
	if err != nil {
		panic(err)
	}

	DeleteAgreement, err = DB.Prepare("DELETE FROM agreement WHERE id = $1")
	if err != nil {
		panic(err)
	}

	UpsertStatus, err = DB.Prepare("SELECT upsert_status($1, $2)")
	if err != nil {
		panic(err)
	}

	GetStatusHistory, err = DB.Prepare("SELECT data FROM agreement WHERE data->>'agreementID' = $1")
	if err != nil {
		panic(err)
	}
}
