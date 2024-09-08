package handler

import (
	"trello-backend/lib"
	"trello-backend/models"
	"trello-backend/pkg/service"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	services service.Settings
}

func NewSettingsHandler(services service.Settings) *SettingsHandler {
	return &SettingsHandler{services: services}
}

func (s *SettingsHandler) GetSettings(c *gin.Context) *lib.WebError {
	payload := c.MustGet("payload")

	customPayload, ok := payload.(*service.Payload)

	if !ok {
		return lib.ServerError("payload is invalid")
	}

	settings, err := s.services.GetSettings(customPayload.UserId)

	if err != nil {
		return err
	}

	c.JSON(200, gin.H{
		"message":  "you've successfully gotten settings",
		"settings": settings,
	})

	return nil
}

type UpdateSettingsInput struct {
	Settings *models.Settings `json:"settings"`
}

func (s *SettingsHandler) UpdateSettings(c *gin.Context) *lib.WebError {
	payload := c.MustGet("payload")

	_, ok := payload.(*service.Payload)

	if !ok {
		return lib.ServerError("payload is invalid")
	}

	var input *UpdateSettingsInput

	if err := c.BindJSON(&input); err != nil {
		return lib.BadRequest("Невалидный запрос")
	}

	if err := s.services.UpdateSettings(*input.Settings); err != nil {
		return err
	}

	c.JSON(200, gin.H{
		"message": "you successfully updated settings",
	})

	return nil
}
