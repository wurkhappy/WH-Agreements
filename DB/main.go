package DB

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	"log"
)

var DB *sql.DB
var Name string = "wurkhappy"

func Setup(production bool) {
	Connect(production)
	CreateStatements()
}

func Connect(production bool) {
	var err error
	if production {
		DB, err = sql.Open("postgres", "user=wurkhappy password=whcollab dbname="+Name+" sslmode=disable")
	} else {
		DB, err = sql.Open("postgres", "user=postgres dbname="+Name+" sslmode=disable")
	}
	if err != nil {
		panic(err)
	}
}

func Close() {
	log.Println("close db")
	SaveAgreement.Close()
	UpsertAgreement.Close()
	FindAgreementByVersionID.Close()
	FindLiveAgreementsByClientID.Close()
	FindAgreementByUserID.Close()
	FindAgreementByFreelancerID.Close()
	FindArchivedByFreelancerID.Close()
	FindArchivedByClientID.Close()
	DeleteAgreement.Close()
	FindLiveVersions.Close()
	UpsertStatus.Close()
	GetStatusHistory.Close()
	DB.Close()
}
