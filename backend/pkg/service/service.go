package service

import (
	"time"
	"trello-backend/lib"
	"trello-backend/models"
	"trello-backend/pkg/repository"

	"github.com/dgrijalva/jwt-go"
)

type AuthServiceReturn struct {
	AccessToken  *TokenServiceReturn
	RefreshToken *TokenServiceReturn
}

type TokenServiceReturn struct {
	Token      string
	Expiration time.Duration
}

type Auth interface {
	SignUp(input *models.AuthInput, currentRefreshToken string, fingerPrint string) (*AuthServiceReturn, *lib.WebError)
	SignIn(input *models.AuthInput, currentRefreshToken string, fingerPrint string) (*AuthServiceReturn, *lib.WebError)
	LogOut(refreshToken string) *lib.WebError
	Refresh(currentRefreshToken string, fingerPrint string) (*AuthServiceReturn, *lib.WebError)
	ChangeUserData(input *models.AuthInput, currentRefreshToken string, fingerPrint string) (*AuthServiceReturn, *lib.WebError)
	GetUserData(refreshToken string) (*models.User, *lib.WebError)
}

type Token interface {
	GenerateAccessToken(claims *Payload) (*TokenServiceReturn, *lib.WebError)
	GenerateRefreshToken(claims *Payload) (*TokenServiceReturn, *lib.WebError)
	ParseRefreshToken(refreshToken string) (*Payload, *lib.WebError)
	ParseAccessToken(accessToken string) (*Payload, *lib.WebError)
}

type TimeBlock interface {
	GetTimeBlocks(id int) ([]models.TimeBlock, *lib.WebError)
	CreateTimeBlock(input *models.TimeBlockInput) (int, *lib.WebError)
	DeleteTimeBlock(userId, id int) *lib.WebError
	UpdateTimeBlock(input *models.TimeBlock) *lib.WebError
	ChangeOrder(input *repository.ChangeOrderProps, userId int) *lib.WebError
}

type Settings interface {
	GetSettings(userId int) (*models.Settings, *lib.WebError)
	CreateSettings(userId int) *lib.WebError
	UpdateSettings(settings models.Settings) *lib.WebError
}

type Service struct {
	Auth
	Token
	TimeBlock
	Settings
}

type Payload struct {
	UserId int
	Name   string
	Role   int8
}

type TokenClaims struct {
	jwt.StandardClaims
	Payload *Payload
}

func NewService(repo *repository.Repository) *Service {
	tokenService := NewTokenService()

	return &Service{
		Auth:      NewAuthService(repo.Auth, *tokenService, repo.Settings),
		Token:     tokenService,
		TimeBlock: NewTimeBlockService(repo.TimeBlock),
		Settings:  NewSettingsService(repo.Settings),
	}
}
