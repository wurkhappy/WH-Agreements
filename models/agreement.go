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
	VersionID        string        `json:"versionID"` //tracks agreements across versions
	Version          int           `json:"version"`
	ClientID         string        `json:"clientID"`
	FreelancerID     string        `json:"freelancerID"`
	Title            string        `json:"title"`
	ProposedServices string        `json:"proposedServices"`
	RefundPolicy     string        `json:"refundPolicy"`
	Payments         []*Payment    `json:"payments"`
	StatusHistory    statusHistory `json:"statusHistory"`
	LastModified     time.Time     `json:"lastModified"`
	Archived         bool          `json:"archived"`
	Final            bool          `json:"final"`
	Draft            bool          `json:"draft"`
	CurrentStatus    *Status       `json:"currentStatus"`
}

func NewAgreement() *Agreement {
	id, _ := uuid.NewV4()
	return &Agreement{
		VersionID:     id.String(),
		StatusHistory: nil,
		Version:       0,
		AgreementID:   id.String(),
		Draft:         true,
	}
}

func (a *Agreement) Save() (err error) {
	a.LastModified = time.Now()
	a.StatusHistory = nil

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

func FindArchivedByFreelancerID(id string) (agrmnts []*Agreement, err error) {
	r, err := DB.FindArchivedByFreelancerID.Query(id)
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

func FindArchivedByClientID(id string) (agrmnts []*Agreement, err error) {
	r, err := DB.FindArchivedByClientID.Query(id)
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

	for r.Next() {
		var s string
		err = r.Scan(&s)
		if err != nil {
			return err
		}
		var agr *Agreement
		json.Unmarshal([]byte(s), &agr)
		agreements = append(agreements, agr)
	}

	for _, ag := range agreements {
		if ag.VersionID != agreement.VersionID {
			ag.Archive()
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

func (a *Agreement) GetPayment(id string) *Payment {
	for _, payment := range a.Payments {
		if payment.ID == id {
			return payment
		}
	}
	return nil
}
