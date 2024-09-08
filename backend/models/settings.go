package models

import "time"

type Settings struct {
	Id       int        `json:"id"`
	UserId   int        `json:"user_id"`
	BirthDay *time.Time `json:"birth_day"`
}
