package models

import (
	"time"
)

type Comment struct {
	AuthorID    string    `json:"authorID"`
	DateCreated time.Time `json:"dateCreated"`
	Text        string    `json:"text"`
}
