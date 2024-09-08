package models

type RefreshSession struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	FingerPrint  string `json:"finger_print"`
}
