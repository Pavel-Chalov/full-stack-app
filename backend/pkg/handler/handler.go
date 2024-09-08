package handler

import (
	"trello-backend/lib"
	"trello-backend/pkg/middlewares"
	"trello-backend/pkg/service"

	"github.com/gin-gonic/gin"
)

const (
	RefreshTokenCookieName = "refresh_token"
	AccessTokenCookieName  = "access_token"
)

type Auth interface {
	SignUp(c *gin.Context) *lib.WebError
	SignIn(c *gin.Context) *lib.WebError
	LogOut(c *gin.Context) *lib.WebError
	Refresh(c *gin.Context) *lib.WebError
	ChangeUserData(c *gin.Context) *lib.WebError
	GetUserData(c *gin.Context) *lib.WebError
}

type TimeBlocks interface {
	GetTimeBlocks(c *gin.Context) *lib.WebError
	CreateTimeBlock(c *gin.Context) *lib.WebError
	DeleteTimeBlock(c *gin.Context) *lib.WebError
	UpdateTimeBlock(c *gin.Context) *lib.WebError
	ChangeOrder(c *gin.Context) *lib.WebError
}

type Settings interface {
	GetSettings(c *gin.Context) *lib.WebError
	UpdateSettings(c *gin.Context) *lib.WebError
}

type Middlewares interface {
	CheckAuth() gin.HandlerFunc
}

type Handler struct {
	Auth
	TimeBlocks
	Settings
	Middlewares
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Auth:        NewAuthHandler(services.Auth),
		TimeBlocks:  NewTimeBlocksHandler(services.TimeBlock),
		Middlewares: middlewares.NewMiddlewares(services),
		Settings:    NewSettingsHandler(services.Settings),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode("release")

	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Укажите здесь ваш домен
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Если это предварительный запрос, то нет необходимости продолжать обработку
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", customHandler(h.Auth.SignUp))
		auth.POST("/sign-in", customHandler(h.Auth.SignIn))
		auth.POST("/log-out", customHandler(h.Auth.LogOut))
		auth.GET("/refresh", customHandler(h.Auth.Refresh))
		auth.PUT("/change-data", customHandler(h.Auth.ChangeUserData))
		auth.GET("/get-data", customHandler(h.Auth.GetUserData))
	}

	resource := router.Group("/resource")
	{
		resource.Use(h.Middlewares.CheckAuth())

		timeBlocks := resource.Group("/time-blocks")
		{
			timeBlocks.GET("/get", customHandler(h.TimeBlocks.GetTimeBlocks))
			timeBlocks.POST("/create", customHandler(h.TimeBlocks.CreateTimeBlock))
			timeBlocks.DELETE("/delete", customHandler(h.TimeBlocks.DeleteTimeBlock))
			timeBlocks.PUT("/update", customHandler(h.TimeBlocks.UpdateTimeBlock))
			timeBlocks.PUT("/reorder", customHandler(h.TimeBlocks.ChangeOrder))
		}
	}

	settings := router.Group("/settings")
	{
		settings.Use(h.Middlewares.CheckAuth())

		settings.GET("/get", customHandler(h.Settings.GetSettings))
		settings.PUT("/update", customHandler(h.Settings.UpdateSettings))
	}

	return router
}

func customHandler(f func(c *gin.Context) *lib.WebError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := f(c); err != nil {
			c.JSON(err.Status, map[string]interface{}{"message": err.Error})
		}
	}
}
