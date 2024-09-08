package models

import "time"

type TimeBlock struct {
	Id       int           `json:"id"`
	UserId   int           `json:"user_id"`
	Name     string        `json:"name"`
	Color    string        `json:"color"`
	Order    int           `json:"ordering"`
	Duration time.Duration `json:"duration"`
}
