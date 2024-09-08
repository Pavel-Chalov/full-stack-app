package service

import (
	"fmt"
	"time"
	"trello-backend/lib"
	"trello-backend/models"
	"trello-backend/pkg/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo         repository.Auth
	settingsRepo repository.Settings
	tokenService TokenService
}

func NewAuthService(repo repository.Auth, tokenService TokenService, settingsRepo repository.Settings) *AuthService {
	return &AuthService{repo: repo, tokenService: tokenService, settingsRepo: settingsRepo}
}

func (s *AuthService) SignUp(input *models.AuthInput, currentRefreshToken string, fingerPrint string) (*AuthServiceReturn, *lib.WebError) {
	if userWithSuchName, err := s.repo.GetUser(input.Name); userWithSuchName != nil && err == nil {
		return nil, lib.Conflict("Пользователь с таким именем уже существует")
	} else if err.Status != 404 && userWithSuchName == nil {
		return nil, err
	}

	hashedPassword, bcryptErr := bcrypt.GenerateFromPassword([]byte(input.Password), 9)

	if bcryptErr != nil {
		return nil, lib.ServerError(bcryptErr.Error())
	}

	input.Password = string(hashedPassword)

	user, err := s.repo.CreateUser(input)

	if err != nil {
		return nil, err
	}

	payload := &Payload{UserId: user.Id, Name: input.Name, Role: user.Role}

	accessToken, err := s.tokenService.GenerateAccessToken(payload)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(payload)

	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteRefreshSession(currentRefreshToken); err != nil {
		return nil, err
	}

	if err := s.repo.CreateRefreshSession(user.Id, refreshToken.Token, fingerPrint); err != nil {
		return nil, err
	}

	if err := s.settingsRepo.CreateSettings(user.Id); err != nil {
		return nil, err
	}

	return &AuthServiceReturn{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) SignIn(input *models.AuthInput, currentRefreshToken string, fingerPrint string) (*AuthServiceReturn, *lib.WebError) {
	user, err := s.repo.GetUser(input.Name)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, lib.Forbidden("Некорректный пароль")
	}

	payload := &Payload{UserId: user.Id, Name: input.Name, Role: user.Role}

	accessToken, err := s.tokenService.GenerateAccessToken(payload)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(payload)

	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteRefreshSession(currentRefreshToken); err != nil {
		return nil, err
	}

	if err := s.repo.CreateRefreshSession(user.Id, refreshToken.Token, fingerPrint); err != nil {
		return nil, err
	}

	return &AuthServiceReturn{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) LogOut(refreshToken string) *lib.WebError {
	if refreshToken == "" {
		return nil
	}

	if err := s.repo.DeleteRefreshSession(refreshToken); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Refresh(currentRefreshToken string, fingerPrint string) (*AuthServiceReturn, *lib.WebError) {
	if currentRefreshToken == "" {
		return nil, lib.Unauthorized("Вы не авторизованы")
	}

	payload, err := s.tokenService.ParseRefreshToken(currentRefreshToken)

	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUser(payload.Name)

	if err != nil {
		return nil, err
	}

	refreshSession, err := s.repo.GetRefreshSession(currentRefreshToken)

	if err != nil {
		if err.Status == 404 {
			return nil, lib.Unauthorized("Вы не авторизованы")
		}

		return nil, err
	}

	if refreshSession.UserId != payload.UserId && refreshSession.UserId != user.Id || refreshSession.FingerPrint != fingerPrint {
		return nil, lib.Unauthorized("Вы не авторизованы")
	}

	newPayload := &Payload{UserId: user.Id, Name: user.Name, Role: user.Role}

	accessToken, err := s.tokenService.GenerateAccessToken(newPayload)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(newPayload)

	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteRefreshSession(currentRefreshToken); err != nil {
		return nil, err
	}

	if err := s.repo.CreateRefreshSession(payload.UserId, refreshToken.Token, fingerPrint); err != nil {
		return nil, err
	}

	return &AuthServiceReturn{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) ChangeUserData(input *models.AuthInput, currentRefreshToken string, fingerPrint string) (*AuthServiceReturn, *lib.WebError) {
	if currentRefreshToken == "" {
		return nil, lib.Unauthorized("Вы не авторизованы")
	}

	payload, err := s.tokenService.ParseRefreshToken(currentRefreshToken)

	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUser(payload.Name)

	if err != nil {
		return nil, err
	}

	if time.Since(user.UpdatedAt) < 24*time.Hour {
		return nil, lib.Conflict(fmt.Sprintf("Вы можите поменять данные %s", 24*time.Hour-time.Since(user.UpdatedAt)))
	}

	_, err = s.repo.GetRefreshSession(currentRefreshToken)

	if err != nil {
		if err.Status == 404 {
			return nil, lib.Unauthorized("Вы не авторизованы")
		}

		return nil, err
	}

	user, err = s.repo.GetUser(input.Name)

	if err != nil {
		if err.Status != 404 {
			return nil, err
		}
	}

	if user != nil {
		if user.Id != payload.UserId {
			return nil, lib.Conflict("Пользователь с таким именем уже существует")
		}
	}

	hashedPassword, bcryptErr := bcrypt.GenerateFromPassword([]byte(input.Password), 9)

	if bcryptErr != nil {
		return nil, lib.ServerError(bcryptErr.Error())
	}

	input.Password = string(hashedPassword)

	if err = s.repo.ChangeUserData(input, payload.UserId); err != nil {
		return nil, err
	}

	payload = &Payload{UserId: payload.UserId, Name: input.Name, Role: payload.Role}

	accessToken, err := s.tokenService.GenerateAccessToken(payload)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(payload)

	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteRefreshSession(currentRefreshToken); err != nil {
		return nil, err
	}

	if err := s.repo.CreateRefreshSession(payload.UserId, refreshToken.Token, fingerPrint); err != nil {
		return nil, err
	}

	return &AuthServiceReturn{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetUserData(accessToken string) (*models.User, *lib.WebError) {
	if accessToken == "" {
		return nil, lib.Unauthorized("Вы не авторизованы")
	}

	payload, err := s.tokenService.ParseAccessToken(accessToken)

	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUser(payload.Name)

	if err != nil {
		return nil, err
	}

	return user, nil
}
