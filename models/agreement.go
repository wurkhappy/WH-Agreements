package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Agreements/DB"
	"labix.org/v2/mgo/bson"
	"time"
)

type Agreement struct {
	ID               string        `json:"id" bson:"_id"`
	AgreementID      string        `json:"agreementID"`
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
}

func NewAgreement() *Agreement {
	id, _ := uuid.NewV4()
	return &Agreement{
		AgreementID: id.String(),
		Version:     1,
		ID:          id.String(),
	}
}

func (a *Agreement) SaveAgreementWithCtx(ctx *DB.Context) (err error) {
	a.LastModified = time.Now()

	a.addIDtoPayments()

	coll := ctx.Database.C("agreements")
	if _, err := coll.UpsertId(a.ID, &a); err != nil {
		return err
	}
	return nil
}

func (a *Agreement) addIDtoPayments() {
	for _, payment := range a.Payments {
		if payment.ID == "" {
			id, _ := uuid.NewV4()
			payment.ID = id.String()
		}
	}
}

func (a *Agreement) GetID() (id string) {
	return a.AgreementID
}

func (a *Agreement) SetClientID(id string) {
	a.ClientID = id
}

func FindAgreementByID(id string, ctx *DB.Context) (a *Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"_id": id, "archived": false}).One(&a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func FindAgreementByClientID(id string, ctx *DB.Context) (agrmnts []*Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"clientid": id, "archived": false}).Sort("-lastmodified").All(&agrmnts)
	if err != nil {
		return nil, err
	}

	return agrmnts, nil
}

func FindAgreementByFreelancerID(id string, ctx *DB.Context) (agrmnts []*Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"freelancerid": id, "archived": false}).Sort("-lastmodified").All(&agrmnts)
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
