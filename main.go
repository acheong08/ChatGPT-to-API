package main

import (
	"freechatgpt/internal/chatgpt"
	"freechatgpt/internal/tokens"
	typings "freechatgpt/internal/typings"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/r3labs/sse/v2"
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
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	/// Admin routes
	router.PATCH("/admin/puid", admin_check, func(c *gin.Context) {
		// Get the puid from the request (json) and update the puid
		type puid_struct struct {
			PUID string `json:"puid"`
		}
		var puid puid_struct
		err := c.BindJSON(&puid)
		if err != nil {
			c.String(400, "puid not provided")
			return
		}
		PUID = puid.PUID
		c.String(200, "puid updated")
	})
	router.PATCH("/admin/password", admin_check, func(c *gin.Context) {
		// Get the password from the request (json) and update the password
		type password_struct struct {
			Password string `json:"password"`
		}
		var password password_struct
		err := c.BindJSON(&password)
		if err != nil {
			c.String(400, "password not provided")
			return
		}
		ADMIN_PASSWORD = password.Password
		c.String(200, "password updated")
	})
	router.PATCH("/admin/tokens", admin_check, func(c *gin.Context) {
		// Get the request_tokens from the request (json) and update the request_tokens
		var request_tokens []string
		err := c.BindJSON(&request_tokens)
		if err != nil {
			c.String(400, "tokens not provided")
			return
		}
		ACCESS_TOKENS = tokens.NewAccessToken(request_tokens)
		c.String(200, "tokens updated")
	})
	/// Public routes
	router.POST("/v1/chat/completions", func(c *gin.Context) {
		var chat_request typings.APIRequest
		err := c.BindJSON(&chat_request)
		if err != nil {
			c.String(400, "chat request not provided")
			return
		}
		// Convert the chat request to a ChatGPT request
		chatgpt_request := chatgpt.ConvertAPIRequest(chat_request)
		// c.JSON(200, chatgpt_request)
		response, err := chatgpt.SendRequest(chatgpt_request, &PUID, ACCESS_TOKENS.GetToken())
		if err != nil {
			c.String(500, "error sending request")
			return
		}
		defer response.Body.Close()
		c.Status(response.StatusCode)
		c.Stream(func(w io.Writer) bool {
			// Write data to client
			io.Copy(w, response.Body)
			return false
		})

	})
	router.Run(HOST + ":" + PORT)
}
