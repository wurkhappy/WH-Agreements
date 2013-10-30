package models

type Clause struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	UserID string `json:"userID"`
}
