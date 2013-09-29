package models

import (
	"github.com/wurkhappy/WH-Agreements/DB"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type statusHistory []*Status

type Status struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	AgreementID bson.ObjectId `json:"agreementID"`
	PaymentID   bson.ObjectId `json:"paymentID" bson:",omitempty"`
	Action      string        `json:"action"`
	Date        time.Time     `json:"date"`
	Comments    []*Comment    `json:"comments"`
}

func StatusCreated(agrmntID string, paymentID string) *Status {
	agrmntid := bson.ObjectIdHex(agrmntID)
	var paymentid bson.ObjectId
	if paymentID != "" {
	 	paymentid = bson.ObjectIdHex(paymentID)
	 } 

	return &Status{ID: bson.NewObjectId(), Date: time.Now(), Action: "created", AgreementID: agrmntid, PaymentID: paymentid}
}
func StatusAccepted(agrmntID string, paymentID string) *Status {
	agrmntid := bson.ObjectIdHex(agrmntID)
	var paymentid bson.ObjectId
	if paymentID != "" {
	 	paymentid = bson.ObjectIdHex(paymentID)
	 } 

	return &Status{ID: bson.NewObjectId(), Date: time.Now(), Action: "accepted", AgreementID: agrmntid, PaymentID: paymentid}
}
func StatusRejected(agrmntID string, paymentID string) *Status {
	agrmntid := bson.ObjectIdHex(agrmntID)
	var paymentid bson.ObjectId
	if paymentID != "" {
	 	paymentid = bson.ObjectIdHex(paymentID)
	 } 

	return &Status{ID: bson.NewObjectId(), Date: time.Now(), Action: "rejected", AgreementID: agrmntid, PaymentID: paymentid}
}
func StatusSubmitted(agrmntID string, paymentID string) *Status {
	agrmntid := bson.ObjectIdHex(agrmntID)
	var paymentid bson.ObjectId
	if paymentID != "" {
	 	paymentid = bson.ObjectIdHex(paymentID)
	 } 

	return &Status{ID: bson.NewObjectId(), Date: time.Now(), Action: "submitted", AgreementID: agrmntid, PaymentID: paymentid}
}

func (s *Status) UpdateAgreementStatus(ctx *DB.Context) (err error) {

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"statushistory.$": &s}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("agreements")
	info, err := coll.Find(bson.M{
		"_id":               s.AgreementID,
		"statushistory._id": s.ID,
	}).Apply(change, nil)

	log.Print(info)
	log.Print(err)

	return err
}

func (s *Status) UpdatePaymentStatus(ctx *DB.Context) (err error) {

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"statushistory.$": &s}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("agreements")
	info, err := coll.Find(bson.M{
		"_id":                          s.AgreementID,
		"payments._id":                 s.PaymentID,
		"payments.$.statushistory._id": s.ID,
	}).Apply(change, nil)

	log.Print(info)
	log.Print(err)

	return err
}

func (s *Status) AddAgreementStatus(ctx *DB.Context) (err error) {
	m := make(map[string]interface{})
	change := mgo.Change{
		Update:    bson.M{"$push": bson.M{"statushistory": &s}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("agreements")
	info, err := coll.Find(bson.M{
		"_id": s.AgreementID,
	}).Apply(change, &m)

	log.Print(info)
	log.Print(err)

	return err
}

func (s *Status) AddPaymentStatus(ctx *DB.Context) (err error) {
	m := make(map[string]interface{})

	change := mgo.Change{
		Update:    bson.M{"$push": bson.M{"payments.$.statushistory": &s}},
		ReturnNew: true,
	}
	coll := ctx.Database.C("agreements")
	info, err := coll.Find(bson.M{
		"_id":          s.AgreementID,
		"payments._id": s.PaymentID,
	}).Apply(change, &m)

	log.Print(info)
	log.Print(err)

	return err
}
