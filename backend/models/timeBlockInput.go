package models

import "time"

type TimeBlockInput struct {
	UserId   int           `json:"user_id"`
	Name     string        `json:"name" binding:"required,max=24"`
	Color    string        `json:"color" binding:"required"`
	Order    int           `json:"ordering" binding:"required"`
	Duration time.Duration `json:"duration" binding:"required"`
}
