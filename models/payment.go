package models

import (
	"labix.org/v2/mgo/bson"
)

type Payment struct {
	ID     bson.ObjectId `bson:"_id"`
	Amount float64
	Items  []*string
}
