package models

import "time"

type Match struct {
	ID       string    `json:"id"`
	Score    string    `json:"score"`
	Overs    string    `json:"overs"`
	Status   string    `json:"status"`
	UpdateAt time.Time `json:"created_at"`
}
