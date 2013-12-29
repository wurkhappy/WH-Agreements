package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Agreements/DB"
	// "labix.org/v2/mgo"
	"encoding/json"
	// "labix.org/v2/mgo/bson"
	"log"
	"time"
)

type statusHistory []*Status

type Status struct {
	ID                 string    `json:"id" bson:"_id"`
	AgreementID        string    `json:"agreementID"`
	AgreementVersionID string    `json:"agreementVersionID"`
	AgreementVersion   int       `json:"agreementVersion"`
	ParentID           string    `json:"paymentID"`
	Action             string    `json:"action"`
	Date               time.Time `json:"date"`
	UserID             string    `json:"userID"`
}

func CreateStatus(agrmntID, versionID, parentID, action string, versionNumber int) *Status {
	id, _ := uuid.NewV4()
	return &Status{
		ID:                 id.String(),
		Date:               time.Now(),
		Action:             action,
		AgreementID:        agrmntID,
		ParentID:           parentID,
		AgreementVersionID: versionID,
		AgreementVersion:   versionNumber,
	}
}

func (s *Status) Save() (err error) {
	jsonByte, _ := json.Marshal(s)
	_, err = DB.UpsertStatus.Query(s.ID, string(jsonByte))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func GetStatusHistory(agreementID string) (statuses []*Status, err error) {
	r, err := DB.GetStatusHistory.Query(agreementID)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for r.Next() {
		var s string
		err = r.Scan(&s)
		if err != nil {
			return nil, err
		}
		var st *Status
		json.Unmarshal([]byte(s), &st)
		statuses = append(statuses, st)
	}
	return statuses, nil
}
