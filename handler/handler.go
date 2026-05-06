package handler

import (
	"discord-server-go/handler/middleware"
	"discord-server-go/model"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService model.UserService
}

type Config struct {
	R           *gin.Engine
	UserService model.UserService
}

func NewHandler(c *Config) {
	h := &Handler{
		userService: c.UserService,
	}

	c.R.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "No route found",
		})
	})
	c.R.Use(static.Serve("/static", static.LocalFile("./static", true)))
	// if gin.Mode() != gin.TestMode {
	// 	c.R.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
	// }

	ag := c.R.Group("/api/account")

	ag.POST("/register", h.Register)
	ag.POST("/login", h.Login)
	ag.POST("/logout", h.Logout)

	ag.Use(middleware.AuthUser())
	ag.GET("", h.GetCurrent)
}

func toFieldErrorResponse(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": []model.FieldError{
			{
				Field:   field,
				Message: message,
			},
		},
	})
}

func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		log.Printf("Error setting the session: %v\n", err.Error())
	}
}
