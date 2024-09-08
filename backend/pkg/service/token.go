package service

import (
	"time"
	"trello-backend/lib"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

const (
	AccessTokenExpiration  = time.Minute * 30
	RefreshTokenExpiration = time.Hour * 24 * 15
)

func (s *TokenService) GenerateAccessToken(payload *Payload) (*TokenServiceReturn, *lib.WebError) {
	claims := &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpiration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(viper.GetString("access_token_secret")))

	if err != nil {
		return nil, lib.ServerError(err.Error())
	}

	return &TokenServiceReturn{Token: tokenString, Expiration: AccessTokenExpiration}, nil
}

func (s *TokenService) GenerateRefreshToken(payload *Payload) (*TokenServiceReturn, *lib.WebError) {
	claims := &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(RefreshTokenExpiration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(viper.GetString("refresh_token_secret")))

	if err != nil {
		return nil, lib.ServerError(err.Error())
	}

	return &TokenServiceReturn{Token: tokenString, Expiration: RefreshTokenExpiration}, nil
}

func (s *TokenService) ParseRefreshToken(refreshToken string) (*Payload, *lib.WebError) {
	token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("refresh_token_secret")), nil
	})

	if err != nil {
		return nil, lib.Unauthorized(err.Error())
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, lib.Unauthorized("Некорректный токен")
	}

	return claims.Payload, nil
}

func (s *TokenService) ParseAccessToken(accessToken string) (*Payload, *lib.WebError) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("access_token_secret")), nil
	})

	if err != nil {
		return nil, lib.Unauthorized(err.Error())
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, lib.Unauthorized("Некорректный токен")
	}

	return claims.Payload, nil
}
