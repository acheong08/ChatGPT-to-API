package main

import (
	"bufio"
	"freechatgpt/internal/chatgpt"
	typings "freechatgpt/internal/typings"
	"freechatgpt/internal/typings/responses"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func passwordHandler(c *gin.Context) {
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
	// Set environment variable
	os.Setenv("ADMIN_PASSWORD", ADMIN_PASSWORD)
	c.String(200, "password updated")
}

func tokensHandler(c *gin.Context) {
	// Get the request_tokens from the request (json) and update the request_tokens
	type auth struct {
		AuthCookie string `json:"auth_cookie"`
	}
	var auth_req auth
	err := c.BindJSON(&auth_req)
	if err != nil {
		c.String(400, "tokens not provided")
		return
	}
	auth_cookie = auth_req.AuthCookie
	c.String(200, "cookies updated")
}
func optionsHandler(c *gin.Context) {
	// Set headers for CORS
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
func nightmare(c *gin.Context) {
	var original_request typings.APIRequest
	err := c.BindJSON(&original_request)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}
	// Convert the chat request to a ChatGPT request
	translated_request := chatgpt.ConvertAPIRequest(original_request)

	response, err := chatgpt.SendRequest(translated_request)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error sending request",
		})
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, gin.H{
			"error": "error sending request", "details": bufio.NewReader(response.Body),
		})
		return
	}
	// Get response body
	fulltext, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error reading response",
		})
		return
	}

	if !original_request.Stream {
		full_response := responses.NewChatCompletion(string(fulltext))
		if err != nil {
			return
		}
		c.JSON(200, full_response)
		return
	}
	c.String(200, "data: [DONE]")

}
