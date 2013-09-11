package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"github.com/wurkhappy/WH-UserService/DB"
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

func NewAgreement() *User {
	return &User{
		DateCreated:  time.Now(),
		LastModified: time.Now(),
		ID:           bson.NewObjectId(),
	}
}

func (u *User) SaveAgreementWithCtx(ctx *DB.Context) (err error) {
	coll := ctx.Database.C("agreements")
	if _, err := coll.UpsertId(u.ID, &u); err != nil {
		return err
	}
	return nil
}

func FindAgreementByID(id interface{}, ctx *DB.Context) (u *User, err error) {
	switch 
	err = ctx.Database.C("agreements").Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func DeleteAgreementWithID(id string, ctx *DB.Context) (err error) {
	err = ctx.Database.C("agreements").RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
