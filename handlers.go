package main

import (
	"bufio"
	"encoding/json"
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

	response, err := chatgpt.SendRequest(translated_request, auth_cookie)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error sending request",
		})
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var error_response map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&error_response)
		if err != nil {
			c.JSON(response.StatusCode, err)
			return
		}
		c.JSON(response.StatusCode, gin.H{
			"error": "error sending request", "details": error_response,
		})
		return
	}
	// Create a bufio.Reader from the response body
	reader := bufio.NewReader(response.Body)

	var fulltext string = ""

	// Read the response byte by byte until a newline character is encountered
	if original_request.Stream {
		// Response content type is text/event-stream
		c.Header("Content-Type", "text/event-stream")
	} else {
		// Response content type is application/json
		c.Header("Content-Type", "application/json")
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		if len(line) < 3 {
			continue
		}
		// Remove the first and last character from the line
		line = line[1 : len(line)-1]

		translated_response := responses.NewChatCompletionChunk(line)

		// Stream the response to the client
		response_string, err := json.Marshal(translated_response)
		if err != nil {
			continue
		}
		if original_request.Stream {
			_, err = c.Writer.WriteString("data: " + string(response_string) + "\n\n")
			if err != nil {
				return
			}
		}

		// Flush the response writer buffer to ensure that the client receives each line as it's written
		c.Writer.Flush()
		fulltext = fulltext + line
	}

	if !original_request.Stream {
		full_response := responses.NewChatCompletion(fulltext)
		if err != nil {
			return
		}
		c.JSON(200, full_response)
		return
	}
	c.String(200, "data: [DONE]")

}
