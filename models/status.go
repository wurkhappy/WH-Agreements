//IMPORTANT INFO:
//Right now this file contains a workaround for statuses and comments.
//MongoDB cannot dynamically query arrays right now more than one level deep.
//This is a problem because when we update comments for payments, we need to update payments.$.statusHistory.$.comments
//Instead payments no longer have a statusHistory, they just hold current status.
//A payment needs to know it's status in order to manage what actions the user can take.
//However, a payment doesn't need to know about it's comments.
//So what I've done here is give an item it's status as well as pass that same status to the agreements status history.
//With this design, a status is only ever one level deep.
//This isn't a terrible design considering that a status on a payment is also a status on an agreement.

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
