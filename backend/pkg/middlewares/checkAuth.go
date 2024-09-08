package middlewares

import (
	"trello-backend/pkg/service"

	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	services *service.Service
}

func NewMiddlewares(services *service.Service) *Middlewares {
	return &Middlewares{services: services}
}

func (m *Middlewares) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, CookieErr := c.Cookie("access_token")

		if CookieErr != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		payload, parseErr := m.services.Token.ParseAccessToken(accessToken)

		if parseErr != nil {
			c.JSON(parseErr.Status, gin.H{
				"message": parseErr.Error,
			})

			c.Abort()
			return
		}

		c.Set("payload", payload)

		c.Next()
	}
}
