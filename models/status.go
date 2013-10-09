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
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type statusHistory []*Status

type Status struct {
	ID          string     `json:"id" bson:"_id"`
	AgreementID string     `json:"agreementID"`
	PaymentID   string     `json:"paymentID" bson:",omitempty"`
	Action      string     `json:"action"`
	Date        time.Time  `json:"date"`
}

func StatusCreated(agrmntID string, paymentID string) *Status {
	id, _ := uuid.NewV4()
	return &Status{ID: id.String(), Date: time.Now(), Action: "created", AgreementID: agrmntID, PaymentID: paymentID}
}
func StatusAccepted(agrmntID string, paymentID string) *Status {
	id, _ := uuid.NewV4()
	return &Status{ID: id.String(), Date: time.Now(), Action: "accepted", AgreementID: agrmntID, PaymentID: paymentID}
}
func StatusRejected(agrmntID string, paymentID string) *Status {
	id, _ := uuid.NewV4()
	return &Status{ID: id.String(), Date: time.Now(), Action: "rejected", AgreementID: agrmntID, PaymentID: paymentID}
}
func StatusSubmitted(agrmntID string, paymentID string) *Status {
	id, _ := uuid.NewV4()
	return &Status{ID: id.String(), Date: time.Now(), Action: "submitted", AgreementID: agrmntID, PaymentID: paymentID}
}

func (s *Status) UpdateAgreementStatus(ctx *DB.Context) (err error) {
	m := make(map[string]interface{})

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"statushistory.$": &s, "lastmodified": time.Now()}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("agreements")
	info, err := coll.Find(bson.M{
		"_id":               s.AgreementID,
		"statushistory._id": s.ID,
	}).Apply(change, &m)

	log.Print(info)
	log.Print(err)

	return err
}

func (s *Status) UpdatePaymentStatus(ctx *DB.Context) (err error) {
	m := make(map[string]interface{})
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"statushistory.$": &s, "lastmodified": time.Now()}},
		ReturnNew: true,
	}

	coll := ctx.Database.C("agreements")
	_, err = coll.Find(bson.M{
		"_id":               s.AgreementID,
		"statushistory._id": s.ID,
	}).Apply(change, &m)

	return err
}

func (s *Status) AddAgreementStatus(ctx *DB.Context) (err error) {
	m := make(map[string]interface{})
	change := mgo.Change{
		Update:    bson.M{"$push": bson.M{"statushistory": &s}, "$set":bson.M{"lastmodified": time.Now(), "currentStatus": &s}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("agreements")
	_, err = coll.Find(bson.M{
		"_id": s.AgreementID,
	}).Apply(change, &m)

	return err
}

func (s *Status) AddPaymentStatus(ctx *DB.Context) (err error) {
	m := make(map[string]interface{})

	change := mgo.Change{
		Update:    bson.M{"$push": bson.M{"statushistory": &s}, "$set": bson.M{"payments.$.currentstatus": &s, "lastmodified": time.Now()}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("agreements")
	_, err = coll.Find(bson.M{
		"_id":          s.AgreementID,
		"payments._id": s.PaymentID,
	}).Apply(change, &m)

	return err
}
