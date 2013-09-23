package models

import (
	"github.com/wurkhappy/WH-Agreements/DB"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

// type agmtPrivateFields struct {
// 	ID               bson.ObjectId `json:"id" bson:"_id"`
// 	ClientID         string        `json:"clientID"`
// 	FreelancerID     string        `json:"freelancerID"`
// 	Title            string        `json:"title"`
// 	ProposedServices string        `json:"proposedServices"`
// 	RefundPolicy     string        `json:"refundPolicy"`
// 	Payments         []*Payment    `json:"payments"`
// 	DateCreated      time.Time
// 	LastModified     time.Time
// 	Status           *status `json:"status`
// }

type Agreement struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	ClientID         string        `json:"clientID"`
	FreelancerID     string        `json:"freelancerID"`
	Title            string        `json:"title"`
	ProposedServices string        `json:"proposedServices"`
	RefundPolicy     string        `json:"refundPolicy"`
	Payments         []*Payment    `json:"payments"`
	DateCreated      time.Time
	LastModified     time.Time
	Status           *status `json:"status"`
}

func NewAgreement() *Agreement {
	return &Agreement{
		DateCreated:  time.Now().UTC(),
		LastModified: time.Now().UTC(),
		ID:           bson.NewObjectId(),
	}
}

func (a *Agreement) SaveAgreementWithCtx(ctx *DB.Context) (err error) {
	a.LastModified = time.Now().UTC()
	coll := ctx.Database.C("agreements")
	if _, err := coll.UpsertId(a.ID, &a); err != nil {
		return err
	}
	return nil
}

func (a *Agreement) GetID() (id bson.ObjectId) {
	return a.ID
}

func (a *Agreement) SetClientID(id string) {
	a.ClientID = id
}

func FindAgreementByID(id interface{}, ctx *DB.Context) (a *Agreement, err error) {
	switch id.(type) {
	case string:
		err = ctx.Database.C("agreements").Find(bson.M{"_id": bson.ObjectIdHex(id.(string))}).One(&a)
	case bson.ObjectId:
		err = ctx.Database.C("agreements").Find(bson.M{"_id": id}).One(&a)
	}
	if err != nil {
		return nil, err
	}

	return a, nil
}

func FindAgreementByClientID(id string, ctx *DB.Context) (agrmnts []*Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"clientid": id}).All(&agrmnts)
	log.Print(err)
	if err != nil {
		return nil, err
	}

	return agrmnts, nil
}

func FindAgreementByFreelancerID(id string, ctx *DB.Context) (agrmnts []*Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"freelancerid": id}).Sort("-datecreated").All(&agrmnts)
	log.Print(err)
	if err != nil {
		return nil, err
	}

	return agrmnts, nil
}

func DeleteAgreementWithID(id string, ctx *DB.Context) (err error) {
	err = ctx.Database.C("agreements").RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
