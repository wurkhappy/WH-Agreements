package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Agreements/DB"
	// "labix.org/v2/mgo"
	// "labix.org/v2/mgo/bson"
	// "log"
	"time"
)

type statusHistory []*Status

type Status struct {
	ID                 string    `json:"id" bson:"_id"`
	AgreementID        string    `json:"agreementID"`
	AgreementVersionID string    `json:"agreementVersionID"`
	AgreementVersion   int       `json:"agreementVersion"`
	PaymentID          string    `json:"paymentID" bson:",omitempty"`
	Action             string    `json:"action"`
	Date               time.Time `json:"date"`
}

func CreateStatus(agrmntID, versionID, paymentID, action string, versionNumber int) *Status {
	id, _ := uuid.NewV4()
	return &Status{
		ID:                 id.String(),
		Date:               time.Now(),
		Action:             action,
		AgreementID:        agrmntID,
		PaymentID:          paymentID,
		AgreementVersionID: versionID,
		AgreementVersion:   versionNumber,
	}
}

func (s *Status) Save(ctx *DB.Context) (err error) {

	coll := ctx.Database.C("status.history")
	if _, err := coll.UpsertId(s.ID, &s); err != nil {
		return err
	}
	return nil
}
