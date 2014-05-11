package models

import (
	"database/sql"
	"encoding/json"
	_ "github.com/bmizerany/pq"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Agreements/DB"
	"log"
	"time"
)

type Agreement struct {
	AgreementID         string    `json:"agreementID,omitempty"`
	VersionID           string    `json:"versionID,omitempty"` //tracks agreements across versions
	Version             int       `json:"version"`
	ClientID            string    `json:"clientID"`
	FreelancerID        string    `json:"freelancerID"`
	Title               string    `json:"title"`
	ProposedServices    string    `json:"proposedServices"`
	LastModified        time.Time `json:"lastModified"`
	LastAction          *Action   `json:"lastAction"`
	LastSubAction       *Action   `json:"lastSubAction"`
	AcceptsCreditCard   bool      `json:"acceptsCreditCard"`
	AcceptsBankTransfer bool      `json:"acceptsBankTransfer"`
	Archived            bool      `json:"archived"`
}

func NewAgreement() *Agreement {
	id, _ := uuid.NewV4()
	return &Agreement{
		VersionID:   id.String(),
		Version:     0, //agreement doesn't get a version number until it has been submitted
		AgreementID: id.String(),
	}
}

func (a *Agreement) Save() (err error) {
	a.LastModified = time.Now()

	jsonByte, _ := json.Marshal(a)
	r, err := DB.UpsertAgreement.Query(a.VersionID, string(jsonByte))
	r.Close()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func FindLatestAgreementByID(id string) (a *Agreement, err error) {
	var s string
	//query sorts by descending and we get the first row so we get the latest
	err = DB.FindLiveVersions.QueryRow(id).Scan(&s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &a)
	return a, nil
}

func FindAgreementByVersionID(id string) (a *Agreement, err error) {
	var s string
	err = DB.FindAgreementByVersionID.QueryRow(id).Scan(&s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &a)
	return a, nil
}
func FindAgreementByVersionNumber(id string, version int) (a *Agreement, err error) {
	var s string
	err = DB.FindAgreementWithVersionNumber.QueryRow(id, version).Scan(&s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &a)
	return a, nil
}

func FindLiveAgreementsByClientID(id string) (agreements []*Agreement, err error) {
	r, err := DB.FindLiveAgreementsByClientID.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return dbRowsToAgreements(r)
}

func FindAgreementByFreelancerID(id string) (agreements []*Agreement, err error) {
	r, err := DB.FindAgreementByFreelancerID.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return dbRowsToAgreements(r)
}

func FindAgreementByUserID(id string) (agreements []*Agreement, err error) {
	r, err := DB.FindAgreementByUserID.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return dbRowsToAgreements(r)
}

func FindArchivedByFreelancerID(id string) (agreements []*Agreement, err error) {
	r, err := DB.FindArchivedByFreelancerID.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return dbRowsToAgreements(r)
}

func FindArchivedByClientID(id string) (agreements []*Agreement, err error) {
	r, err := DB.FindArchivedByClientID.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return dbRowsToAgreements(r)
}

func DeleteAgreementWithVersionID(id string) (err error) {
	_, err = DB.DeleteAgreement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agreement) Archive() {
	a.Archived = true
	err := a.Save()
	if err != nil {
		log.Print(err)
	}
}

func (agreement *Agreement) ArchiveOtherVersions() error {
	var agreements []*Agreement
	r, err := DB.FindLiveVersions.Query(agreement.AgreementID)
	if err != nil {
		return err
	}
	defer r.Close()

	agreements, err = dbRowsToAgreements(r)
	if err != nil {
		return err
	}

	for _, ag := range agreements {
		if ag.VersionID != agreement.VersionID {
			ag.Archive()
		}
	}
	return nil
}

func (a *Agreement) SetRecipient(id string) {
	//assumes that one of these fields (clientID or freelancerID) is set before the recipient
	if a.ClientID == "" {
		a.ClientID = id
	} else if a.FreelancerID == "" {
		a.FreelancerID = id
	}
}

func (a *Agreement) SetAsCompleted() {
	a.LastAction = new(Action)
	a.LastAction.Name = ActionCompleted
	a.LastAction.Date = time.Now()
}

func dbRowsToAgreements(r *sql.Rows) (agreements []*Agreement, err error) {
	for r.Next() {
		var s string
		err = r.Scan(&s)
		if err != nil {
			return nil, err
		}
		var a *Agreement
		json.Unmarshal([]byte(s), &a)
		agreements = append(agreements, a)
	}
	return agreements, nil
}
