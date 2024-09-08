package service

import (
	"time"
	"trello-backend/lib"
	"trello-backend/models"
	"trello-backend/pkg/repository"
)

type TimeBlockService struct {
	repo repository.TimeBlock
}

func NewTimeBlockService(repo repository.TimeBlock) *TimeBlockService {
	return &TimeBlockService{repo: repo}
}

func (s *TimeBlockService) GetTimeBlocks(id int) ([]models.TimeBlock, *lib.WebError) {
	return s.repo.GetTimeBlocks(id)
}

func (s *TimeBlockService) CreateTimeBlock(input *models.TimeBlockInput) (int, *lib.WebError) {
	timeBlocks, err := s.GetTimeBlocks(input.UserId)

	if err != nil {
		return 0, err
	}

	sum := input.Duration

	for i := 0; i < len(timeBlocks); i++ {
		sum += timeBlocks[i].Duration
	}

	if sum > 24*time.Hour {
		return 0, lib.Conflict("duration cannot be more than 24 hours")
	}

	return s.repo.CreateTimeBlock(input)
}

func (s *TimeBlockService) DeleteTimeBlock(userId, id int) *lib.WebError {
	return s.repo.DeleteTimeBlock(userId, id)
}

func (s *TimeBlockService) UpdateTimeBlock(input *models.TimeBlock) *lib.WebError {
	timeBlocks, err := s.repo.GetTimeBlocks(input.UserId)

	if err != nil {
		return err
	}

	var sum time.Duration

	for _, block := range timeBlocks {
		if block.Id == input.Id {
			sum += input.Duration - block.Duration
		}

		sum += block.Duration
	}

	if sum > 24*time.Hour {
		return lib.Conflict("duration cannot be more than 24 hours")
	}

	return s.repo.UpdateTimeBlock(input)
}

func (s *TimeBlockService) ChangeOrder(input *repository.ChangeOrderProps, userId int) *lib.WebError {
	return s.repo.ChangeOrder(input, userId)
}
