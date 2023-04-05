package main

import (
	"freechatgpt/internal/tokens"
	"os"

	"github.com/gin-gonic/gin"
)

var HOST string
var PORT string
var PUID string
var ACCESS_TOKENS tokens.AccessToken

func init() {
	HOST = os.Getenv("SERVER_HOST")
	PORT = os.Getenv("SERVER_PORT")
	PUID = os.Getenv("PUID")
	if HOST == "" {
		HOST = "127.0.0.1"
	}
	if PORT == "" {
		PORT = "8080"
	}
}

func main() {
	router := gin.Default()

	router.Use(cors)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	admin_routes := router.Group("/admin")
	admin_routes.Use(adminCheck)

	/// Admin routes
	admin_routes.PATCH("/puid", puidHandler)
	admin_routes.PATCH("/password", passwordHandler)
	admin_routes.PATCH("/tokens", adminCheck, tokensHandler)
	/// Public routes
	router.OPTIONS("/v1/chat/completions", optionsHandler)
	router.POST("/v1/chat/completions", nightmare)
	router.Run(HOST + ":" + PORT)
}
