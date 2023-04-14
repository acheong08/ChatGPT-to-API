package main

import (
	"bufio"
	"encoding/json"
	"freechatgpt/internal/chatgpt"
	"freechatgpt/internal/tokens"
	typings "freechatgpt/internal/typings"
	"freechatgpt/internal/typings/responses"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func puidHandler(c *gin.Context) {
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
}
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
	// Throw error when model contains gpt-4
	if strings.Contains(original_request.Model, "gpt-4") {
		c.JSON(400, gin.H{
			"error": "gpt-4 is not supported",
		})
		return
	}
	// Convert the chat request to a ChatGPT request
	translated_request := chatgpt.ConvertAPIRequest(original_request)
	// c.JSON(200, chatgpt_request)

	authHeader := c.GetHeader("Authorization")
	token := ACCESS_TOKENS.GetToken()
	if authHeader != "" {
		// 如果Authorization头不为空，则提取其中的token
		// 首先将Bearer前缀替换为空字符串
		customAccessToken := strings.Replace(authHeader, "Bearer ", "", 1)
		if customAccessToken != "" {
			token = customAccessToken
			println("customAccessToken set:" + customAccessToken)
		}
	}

	response, err := chatgpt.SendRequest(translated_request, &PUID, token)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "error sending request",
		})
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, gin.H{
			"error": "error sending request",
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
				println(line)
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
			var delta responses.Delta = responses.Delta{
				Content: original_response.Message.Content.Parts[0],
				Role:    "assistant",
			}
			translated_response := responses.ChatCompletionChunk{
				ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
				Object:  "chat.completion.chunk",
				Created: int64(original_response.Message.CreateTime),
				Model:   "gpt-3.5-turbo-0301",
				Choices: []responses.Choices{
					{
						Index:        0,
						Delta:        delta,
						FinishReason: nil,
					},
				},
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
				full_response := responses.ChatCompletion{
					ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
					Object:  "chat.completion",
					Created: int64(0),
					Model:   "gpt-3.5-turbo-0301",
					Choices: []responses.Choice{
						{
							Message: responses.Msg{
								Content: fulltext,
								Role:    "assistant",
							},
							Index:        0,
							FinishReason: "stop",
						},
					},
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
