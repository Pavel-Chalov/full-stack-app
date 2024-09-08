package models

type AuthInput struct {
	Name     string `json:"name" binding:"required,min=3,max=24"`
	Password string `json:"password" binding:"required,min=8,max=48"`
}
