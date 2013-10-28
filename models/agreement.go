package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Agreements/DB"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
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

func (a *Agreement) SaveAgreementWithCtx(ctx *DB.Context) (err error) {
	a.LastModified = time.Now()

	coll := ctx.Database.C("agreements")
	if _, err := coll.UpsertId(a.VersionID, &a); err != nil {
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

func FindLatestAgreementByID(id string, ctx *DB.Context) (a *Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"agreementid": id}).Sort("-version").One(&a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func FindAgreementByVersionID(id string, ctx *DB.Context) (a *Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"_id": id}).One(&a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func FindLiveAgreementsByClientID(id string, ctx *DB.Context) (agrmnts []*Agreement, err error) {
	err = ctx.Database.C("agreements").Find(bson.M{"clientid": id, "archived": false, "draft": false}).Sort("-lastmodified").All(&agrmnts)
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

func DeleteAgreementWithVersionID(id string, ctx *DB.Context) (err error) {
	err = ctx.Database.C("agreements").RemoveId(id)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agreement) Archive(ctx *DB.Context) {
	m := make(map[string]interface{})

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"archived": true}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("agreements")
	_, _ = coll.Find(bson.M{
		"_id": a.VersionID,
	}).Apply(change, &m)

}

func ArchiveLastAgrmntVersion(id string, ctx *DB.Context) {

	var agreements []*Agreement
	_ = ctx.Database.C("agreements").Find(bson.M{"agreementid": id, "archived": false}).Sort("-version").All(&agreements)
	log.Print(agreements)
	if len(agreements) > 1 {
		agreement := agreements[1]
		agreement.Archive(ctx)

	}
}

func (a *Agreement) SetPaymentStatus(status *Status) {
	for _, payment := range a.Payments {
		if payment.ID == status.PaymentID {
			payment.CurrentStatus = status
		}
	}
}

func (a *Agreement) GetStatusHistory(ctx *DB.Context) []*Status {
	var statusHistory []*Status
	_ = ctx.Database.C("status.history").Find(bson.M{"agreementid": a.AgreementID}).All(&statusHistory)
	return statusHistory
}
