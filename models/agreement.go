package models

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Agreements/DB"
	"log"
	"time"
)

type Agreement struct {
	AgreementID      string        `json:"agreementID"`
	VersionID        string        `json:"versionID" bson:"_id"` //tracks agreements across versions
	Version          int           `json:"version"`
	ClientID         string        `json:"clientID"`
	FreelancerID     string        `json:"freelancerID"`
	Title            string        `json:"title"`
	ProposedServices string        `json:"proposedServices"`
	RefundPolicy     string        `json:"refundPolicy"`
	Payments         []*Payment    `json:"payments"`
	StatusHistory    statusHistory `json:"statusHistory" bson:"-"`
	LastModified     time.Time     `json:"lastModified"`
	Archived         bool          `json:"archived"`
	Draft            bool          `json:"draft"`
	CurrentStatus    *Status       `json:"currentStatus" bson:",omitempty"`
	Clauses          []*Clause     `json:"clauses" bson:",omitempty"`
}

func NewAgreement() *Agreement {
	id, _ := uuid.NewV4()
	return &Agreement{
		VersionID:     id.String(),
		StatusHistory: nil,
		Version:       1,
		AgreementID:   id.String(),
		Draft:         true,
	}
}

func (a *Agreement) Save() (err error) {
	a.LastModified = time.Now()

	jsonByte, _ := json.Marshal(a)
	_, err = DB.UpsertAgreement.Query(a.VersionID, string(jsonByte))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (a *Agreement) AddIDtoPayments() {
	for _, payment := range a.Payments {
		if payment.ID == "" {
			id, _ := uuid.NewV4()
			payment.ID = id.String()
		}
	}
}

func FindLatestAgreementByID(id string) (a *Agreement, err error) {
	var s string
	//query sorts by DESC and we get the first row so we get the latest
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

func FindLiveAgreementsByClientID(id string) (agrmnts []*Agreement, err error) {
	r, err := DB.FindLiveAgreementsByClientID.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for r.Next() {
		var s string
		err = r.Scan(&s)
		if err != nil {
			return nil, err
		}
		var a *Agreement
		json.Unmarshal([]byte(s), &a)
		agrmnts = append(agrmnts, a)
	}
	return agrmnts, nil
}

func FindAgreementByFreelancerID(id string) (agrmnts []*Agreement, err error) {
	r, err := DB.FindAgreementByFreelancerID.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for r.Next() {
		var s string
		err = r.Scan(&s)
		if err != nil {
			return nil, err
		}
		var a *Agreement
		json.Unmarshal([]byte(s), &a)
		agrmnts = append(agrmnts, a)
	}
	return agrmnts, nil
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
	a.Save()
}

func ArchiveLastAgrmntVersion(id string) error {

	var agreements []*Agreement
	r, err := DB.FindLiveVersions.Query(id)
	if err != nil {
		return err
	}
	defer r.Close()

	for r.Next() {
		var s string
		err = r.Scan(&s)
		if err != nil {
			return err
		}
		var a *Agreement
		json.Unmarshal([]byte(s), &a)
		agreements = append(agreements, a)
	}

	count := len(agreements)
	if count > 1 {
		for i := 1; i < count; i++ {
			agreement := agreements[i]
			agreement.Archive()
		}
	}
	return nil
}

func (a *Agreement) SetPaymentStatus(status *Status) {
	for _, payment := range a.Payments {
		if payment.ID == status.PaymentID {
			payment.CurrentStatus = status
		}
	}
}

func (a *Agreement) GetFirstOutstandingPayment() *Payment {
	for _, payment := range a.Payments {
		if payment.CurrentStatus.Action != "accepted" {
			return payment
		}
	}
	return nil
}
