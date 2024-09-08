package handler

import (
	"trello-backend/lib"
	"trello-backend/models"
	"trello-backend/pkg/repository"
	"trello-backend/pkg/service"

	"github.com/gin-gonic/gin"
)

type TimeBlocksHandler struct {
	services service.TimeBlock
}

func NewTimeBlocksHandler(services service.TimeBlock) *TimeBlocksHandler {
	return &TimeBlocksHandler{services: services}
}

func (h *TimeBlocksHandler) GetTimeBlocks(c *gin.Context) *lib.WebError {
	payload := c.MustGet("payload")

	customPayload, ok := payload.(*service.Payload)

	if !ok {
		return lib.ServerError("payload is invalid")
	}

	timeBlocks, err := h.services.GetTimeBlocks(customPayload.UserId)

	if err != nil {
		return err
	}

	c.JSON(200, gin.H{
		"message":    "you've successfully gotten time blocks",
		"timeBlocks": timeBlocks,
	})

	return nil
}

func (h *TimeBlocksHandler) CreateTimeBlock(c *gin.Context) *lib.WebError {
	var input *models.TimeBlockInput

	if err := c.BindJSON(&input); err != nil {
		return lib.BadRequest(err.Error())
	}

	payload, ok := c.MustGet("payload").(*service.Payload)

	if !ok {
		return lib.BadRequest("bad request")
	}

	input.UserId = payload.UserId

	id, err := h.services.CreateTimeBlock(input)

	if err != nil {
		return err
	}

	c.JSON(200, gin.H{
		"message": "success",
		"id":      id,
	})

	return nil
}

type DeleteTimeBlockInput struct {
	Id int `json:"id" binding:"required"`
}

func (h *TimeBlocksHandler) DeleteTimeBlock(c *gin.Context) *lib.WebError {
	payload, ok := c.MustGet("payload").(*service.Payload)

	if !ok {
		return lib.ServerError("payload is invalid")
	}

	var input *DeleteTimeBlockInput

	if err := c.BindJSON(&input); err != nil {
		return lib.BadRequest(err.Error())
	}

	if err := h.services.DeleteTimeBlock(payload.UserId, input.Id); err != nil {
		return err
	}

	c.JSON(200, gin.H{
		"message": "success",
	})

	return nil
}

func (h *TimeBlocksHandler) UpdateTimeBlock(c *gin.Context) *lib.WebError {
	var input *models.TimeBlock

	if err := c.BindJSON(&input); err != nil {
		return lib.BadRequest(err.Error())
	}

	payload, ok := c.MustGet("payload").(*service.Payload)

	if !ok {
		return lib.BadRequest("bad request")
	}

	input.UserId = payload.UserId

	if err := h.services.UpdateTimeBlock(input); err != nil {
		return err
	}

	c.JSON(200, gin.H{
		"message": "success",
	})

	return nil
}

func (h *TimeBlocksHandler) ChangeOrder(c *gin.Context) *lib.WebError {
	var input *repository.ChangeOrderProps

	if err := c.BindJSON(&input); err != nil {
		return lib.BadRequest(err.Error())
	}

	payload, ok := c.MustGet("payload").(*service.Payload)

	if !ok {
		return lib.BadRequest("bad request")
	}

	if err := h.services.ChangeOrder(input, payload.UserId); err != nil {
		return err
	}

	c.JSON(200, gin.H{
		"message": "success",
	})

	return nil
}
