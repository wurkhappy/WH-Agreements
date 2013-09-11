package models

import (
	"github.com/wurkhappy/WH-Agreements/DB"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type Agreement struct {
	ID           bson.ObjectId `bson:"_id"`
	ClientID     bson.ObjectId `bson:"_id"`
	FreelancerID bson.ObjectId `bson:"_id"`
	Title        string
	Description  []string
	Payments     []*Payment
	DateCreated  time.Time
	LastModified time.Time
	Status       *status
}

func NewAgreement() *Agreement {
	return &Agreement{
		DateCreated:  time.Now(),
		LastModified: time.Now(),
		ID:           bson.NewObjectId(),
	}
}

func (a *Agreement) SaveAgreementWithCtx(ctx *DB.Context) (err error) {
	coll := ctx.Database.C("agreements")
	if _, err := coll.UpsertId(a.ID, &a); err != nil {
		return err
	}
	return nil
}

func FindAgreementByID(id interface{}, ctx *DB.Context) (a *Agreement, err error) {
	switch id.(type){
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

func DeleteAgreementWithID(id string, ctx *DB.Context) (err error) {
	err = ctx.Database.C("agreements").RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
