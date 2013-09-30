package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Agreements/DB"
	"labix.org/v2/mgo/bson"
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
	ID               string        `json:"id" bson:"_id"`
	ClientID         string        `json:"clientID"`
	FreelancerID     string        `json:"freelancerID"`
	Title            string        `json:"title"`
	ProposedServices string        `json:"proposedServices"`
	RefundPolicy     string        `json:"refundPolicy"`
	Payments         []*Payment    `json:"payments"`
	StatusHistory    statusHistory `json:"statusHistory"`
	LastModified     time.Time     `json:"-"`
}

func NewAgreement() *Agreement {
	id, _ := uuid.NewV4()
	return &Agreement{
		StatusHistory: statusHistory{
			StatusCreated(id.String(), ""),
		},
		ID: id.String(),
	}
}

func (a *Agreement) SaveAgreementWithCtx(ctx *DB.Context) (err error) {
	a.LastModified = time.Now()

	for _, payment := range a.Payments {
		if payment.ID == "" {
			id, _ := uuid.NewV4()
			payment.ID = id.String()
		}
	}

	coll := ctx.Database.C("agreements")
	if _, err := coll.UpsertId(a.ID, &a); err != nil {
		return err
	}
	return nil
}

func (a *Agreement) GetID() (id string) {
	return a.ID
}

func (a *Agreement) SetClientID(id string) {
	a.ClientID = id
}

func FindAgreementByID(id string, ctx *DB.Context) (a *Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"_id": id}).One(&a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func FindAgreementByClientID(id string, ctx *DB.Context) (agrmnts []*Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"clientid": id}).Sort("-lastmodified").All(&agrmnts)
	if err != nil {
		return nil, err
	}

	return agrmnts, nil
}

func FindAgreementByFreelancerID(id string, ctx *DB.Context) (agrmnts []*Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"freelancerid": id}).Sort("-lastmodified").All(&agrmnts)
	if err != nil {
		return nil, err
	}

	return agrmnts, nil
}

func DeleteAgreementWithID(id string, ctx *DB.Context) (err error) {
	err = ctx.Database.C("agreements").RemoveId(id)
	if err != nil {
		return err
	}
	return nil
}
