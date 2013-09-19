package models

import (
	"encoding/json"
	"github.com/wurkhappy/WH-Agreements/DB"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type agmtPrivateFields struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	ClientID         string        `json:"clientID"`
	FreelancerID     string        `json:"freelancerID"`
	Title            string        `json:"title"`
	ProposedServices string        `json:"proposedServices"`
	RefundPolicy     string        `json:"refundPolicy"`
	Payments         []*Payment    `json:"payments"`
	DateCreated      time.Time
	LastModified     time.Time
	Status           *status `json:"status`
}

type Agreement struct {
	agmtPrivateFields
}

func NewAgreement() *Agreement {
	return &Agreement{
		agmtPrivateFields{
			DateCreated:  time.Now(),
			LastModified: time.Now(),
			ID:           bson.NewObjectId(),
		},
	}
}

func (a *Agreement) SaveAgreementWithCtx(ctx *DB.Context) (err error) {
	a.agmtPrivateFields.LastModified = time.Now()
	coll := ctx.Database.C("agreements")
	if _, err := coll.UpsertId(a.agmtPrivateFields.ID, &a.agmtPrivateFields); err != nil {
		return err
	}
	return nil
}

func (a *Agreement) GetID() (id bson.ObjectId) {
	return a.agmtPrivateFields.ID
}

func (a *Agreement) SetClientID(id string) {
	a.agmtPrivateFields.ClientID = id
}

func (a *Agreement) GetJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":               a.agmtPrivateFields.ID,
		"clientID":         a.agmtPrivateFields.ClientID,
		"freelancerID":     a.agmtPrivateFields.FreelancerID,
		"title":            a.agmtPrivateFields.Title,
		"proposedServices": a.agmtPrivateFields.ProposedServices,
		"refundPolicy":     a.agmtPrivateFields.RefundPolicy,
		"payments":         a.agmtPrivateFields.Payments,
		"status":           a.agmtPrivateFields.Status,
	})
}

func (a *Agreement) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &a.agmtPrivateFields)
}

func FindAgreementByID(id interface{}, ctx *DB.Context) (a *Agreement, err error) {
	a = new(Agreement)
	var fields agmtPrivateFields

	switch id.(type) {
	case string:
		err = ctx.Database.C("agreements").Find(bson.M{"_id": bson.ObjectIdHex(id.(string))}).One(&fields)
	case bson.ObjectId:
		err = ctx.Database.C("agreements").Find(bson.M{"_id": id}).One(&fields)
	}
	if err != nil {
		return nil, err
	}

	a.agmtPrivateFields = fields

	return a, nil
}

func DeleteAgreementWithID(id string, ctx *DB.Context) (err error) {
	err = ctx.Database.C("agreements").RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
