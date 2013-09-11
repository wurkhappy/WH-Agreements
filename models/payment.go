package models

type Payment struct {
	ID     bson.ObjectId `bson:"_id"`
	Amount float64
	Items  []*string
}
