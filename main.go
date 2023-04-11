package main

import (
	"encoding/json"
	"freechatgpt/internal/tokens"
	"os"

	"github.com/fvbock/endless"
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
	// Check if access_tokens.json exists
	if _, err := os.Stat("access_tokens.json"); os.IsNotExist(err) {
		// Create the file
		file, err := os.Create("access_tokens.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()
	} else {
		// Load the tokens
		file, err := os.Open("access_tokens.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		var token_list []string
		err = decoder.Decode(&token_list)
		if err != nil {
			return
		}
		ACCESS_TOKENS = tokens.NewAccessToken(token_list)
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
	endless.ListenAndServe(HOST+":"+PORT, router)
}
