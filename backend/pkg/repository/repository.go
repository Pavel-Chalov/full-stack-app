package repository

import (
	"trello-backend/lib"
	"trello-backend/models"

	"github.com/jmoiron/sqlx"
)

type Auth interface {
	CreateUser(input *models.AuthInput) (*models.User, *lib.WebError)
	GetUser(name string) (*models.User, *lib.WebError)
	ChangeUserData(input *models.AuthInput, id int) *lib.WebError

	GetRefreshSession(refreshToken string) (*models.RefreshSession, *lib.WebError)
	CreateRefreshSession(id int, refreshToken string, fingerPrint string) *lib.WebError
	DeleteRefreshSession(refreshToken string) *lib.WebError
}

type TimeBlock interface {
	GetTimeBlocks(id int) ([]models.TimeBlock, *lib.WebError)
	CreateTimeBlock(input *models.TimeBlockInput) (int, *lib.WebError)
	DeleteTimeBlock(userId, id int) *lib.WebError
	UpdateTimeBlock(input *models.TimeBlock) *lib.WebError
	ChangeOrder(input *ChangeOrderProps, userId int) *lib.WebError
}

type Settings interface {
	GetSettings(userId int) (*models.Settings, *lib.WebError)
	CreateSettings(userId int) *lib.WebError
	UpdateSettings(settings models.Settings) *lib.WebError
}

type Repository struct {
	Auth
	TimeBlock
	Settings
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:      NewAuthPostgres(db),
		TimeBlock: NewTimeBlockPostgres(db),
		Settings:  NewSettingsPostgres(db),
	}
}
