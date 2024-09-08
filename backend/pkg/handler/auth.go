package handler

import (
	"trello-backend/lib"
	"trello-backend/models"
	"trello-backend/pkg/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	services service.Auth
}

func NewAuthHandler(services service.Auth) *AuthHandler {
	return &AuthHandler{services: services}
}

func (h *AuthHandler) SignUp(c *gin.Context) *lib.WebError {
	var input *models.AuthInput

	if err := c.BindJSON(&input); err != nil {
		return lib.BadRequest("Невалидный запрос")
	}

	currentRefreshToken, _ := c.Cookie(RefreshTokenCookieName)

	result, err := h.services.SignUp(input, currentRefreshToken, c.Request.UserAgent())

	if err != nil {
		return err
	}

	c.SetCookie(AccessTokenCookieName, result.AccessToken.Token, int(result.AccessToken.Expiration), "/", "localhost", true, true)
	c.SetCookie(RefreshTokenCookieName, result.RefreshToken.Token, int(result.RefreshToken.Expiration), "/", "localhost", true, true)

	c.JSON(200, map[string]interface{}{
		"message": "Вы успешно зарегистрировались",
	})

	return nil
}

func (h *AuthHandler) SignIn(c *gin.Context) *lib.WebError {
	var input *models.AuthInput

	if err := c.BindJSON(&input); err != nil {
		return lib.BadRequest("Невалидный запрос")
	}

	currentRefreshToken, _ := c.Cookie(RefreshTokenCookieName)

	result, err := h.services.SignIn(input, currentRefreshToken, c.Request.UserAgent())

	if err != nil {
		return err
	}

	c.SetCookie(AccessTokenCookieName, result.AccessToken.Token, int(result.AccessToken.Expiration), "/", "localhost", true, true)
	c.SetCookie(RefreshTokenCookieName, result.RefreshToken.Token, int(result.RefreshToken.Expiration), "/", "localhost", true, true)

	c.JSON(200, map[string]interface{}{
		"message": "Вы успешно вошли",
	})

	return nil
}

func (h *AuthHandler) LogOut(c *gin.Context) *lib.WebError {
	refreshToken, _ := c.Cookie(RefreshTokenCookieName)

	if err := h.services.LogOut(refreshToken); err != nil {
		return err
	}

	c.SetCookie(RefreshTokenCookieName, "", 0, "/", "localhost", true, true)
	c.SetCookie(AccessTokenCookieName, "", 0, "/", "localhost", true, true)

	c.JSON(200, map[string]interface{}{
		"message": "Вы успешно вышли",
	})

	return nil
}

func (h *AuthHandler) Refresh(c *gin.Context) *lib.WebError {
	currentRefreshToken, _ := c.Cookie(RefreshTokenCookieName)

	result, err := h.services.Refresh(currentRefreshToken, c.Request.UserAgent())

	if err != nil {
		return err
	}

	c.SetCookie(AccessTokenCookieName, result.AccessToken.Token, int(result.AccessToken.Expiration), "/", "localhost", true, true)
	c.SetCookie(RefreshTokenCookieName, result.RefreshToken.Token, int(result.RefreshToken.Expiration), "/", "localhost", true, true)

	c.JSON(200, map[string]interface{}{
		"message": "Успех",
	})

	return nil
}

func (h *AuthHandler) ChangeUserData(c *gin.Context) *lib.WebError {
	var input *models.AuthInput

	if err := c.BindJSON(&input); err != nil {
		return lib.BadRequest("Невалидный запрос")
	}

	currentRefreshToken, _ := c.Cookie(RefreshTokenCookieName)

	result, err := h.services.ChangeUserData(input, currentRefreshToken, c.Request.UserAgent())

	if err != nil {
		return err
	}

	c.SetCookie(AccessTokenCookieName, result.AccessToken.Token, int(result.AccessToken.Expiration), "/", "localhost", true, true)
	c.SetCookie(RefreshTokenCookieName, result.RefreshToken.Token, int(result.RefreshToken.Expiration), "/", "localhost", true, true)

	c.JSON(200, map[string]interface{}{
		"message": "Вы успешно изменили свои данные",
	})

	return nil
}

func (h *AuthHandler) GetUserData(c *gin.Context) *lib.WebError {
	accessToken, _ := c.Cookie(AccessTokenCookieName)

	user, err := h.services.GetUserData(accessToken)

	if err != nil {
		return err
	}

	c.JSON(200, map[string]interface{}{
		"user":    user,
		"message": "Успех",
	})

	return nil
}
