package main

import (
	"bufio"
	"encoding/json"
	"freechatgpt/internal/chatgpt"
	"freechatgpt/internal/tokens"
	typings "freechatgpt/internal/typings"
	"freechatgpt/internal/typings/responses"
	"io"
	"os"
	"strings"

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
	var request_tokens []string
	err := c.BindJSON(&request_tokens)
	if err != nil {
		c.String(400, "tokens not provided")
		return
	}
	ACCESS_TOKENS = tokens.NewAccessToken(request_tokens)
	c.String(200, "tokens updated")
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
	// // Throw error when model contains gpt-4
	// if strings.Contains(original_request.Model, "gpt-4") {
	// 	c.JSON(400, gin.H{
	// 		"error": "gpt-4 is not supported",
	// 	})
	// 	return
	// }
	// Convert the chat request to a ChatGPT request
	translated_request := chatgpt.ConvertAPIRequest(original_request)

	if original_request.Model == "gpt-4" {
		translated_request.Model = "gpt-4"
	}

	// authHeader := c.GetHeader("Authorization")
	token := ACCESS_TOKENS.GetToken()
	// if authHeader != "" {
	// 	customAccessToken := strings.Replace(authHeader, "Bearer ", "", 1)
	// 	if customAccessToken != "" {
	// 		token = customAccessToken
	// 		println("customAccessToken set:" + customAccessToken)
	// 	}
	// }

	response, err := chatgpt.SendRequest(translated_request, token)
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

	var fulltext string

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
		if len(line) < 6 {
			continue
		}
		// Remove "data: " from the beginning of the line
		line = line[6:]
		// Check if line starts with [DONE]
		if !strings.HasPrefix(line, "[DONE]") {
			// Parse the line as JSON
			var original_response responses.Data
			err = json.Unmarshal([]byte(line), &original_response)
			if err != nil {
				continue
			}
			if original_response.Error != nil {
				return
			}
			if original_response.Message.Content.Parts[0] == "" || original_response.Message.Author.Role != "assistant" {
				continue
			}
			if original_response.Message.Metadata.Timestamp == "absolute" {
				continue
			}
			tmp_fulltext := original_response.Message.Content.Parts[0]
			original_response.Message.Content.Parts[0] = strings.ReplaceAll(original_response.Message.Content.Parts[0], fulltext, "")
			translated_response := responses.NewChatCompletionChunk(original_response.Message.Content.Parts[0])
			if original_request.Model != "" {
				translated_response.Model = original_request.Model
			}
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
			fulltext = tmp_fulltext
		} else {
			if !original_request.Stream {
				full_response := responses.NewChatCompletion(fulltext)
				if original_request.Model != "" {
					full_response.Model = original_request.Model
				}
				if err != nil {
					return
				}
				c.JSON(200, full_response)
				return
			}
			c.String(200, "data: [DONE]")
			break

		}
	}

}
