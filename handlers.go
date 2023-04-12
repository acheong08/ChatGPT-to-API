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
	var chat_request typings.APIRequest
	err := c.BindJSON(&chat_request)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	// Convert the chat request to a ChatGPT request
	chatgpt_request := chatgpt.ConvertAPIRequest(chat_request)
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

	response, err := chatgpt.SendRequest(chatgpt_request, &PUID, token)
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
	if chat_request.Stream {
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
			c.JSON(500, gin.H{
				"error": "error reading response",
			})
			return
		}
		if len(line) < 6 {
			continue
		}
		// Remove "data: " from the beginning of the line
		line = line[6:]
		// Parse the line as JSON
		var chat_response responses.Data
		err = json.Unmarshal([]byte(line), &chat_response)
		if err != nil {
			continue
		}
		if chat_response.Error != nil {
			c.JSON(500, gin.H{
				"error": chat_response.Error,
			})
			return
		}
		if chat_response.Message.Content.Parts[0] == "" || chat_response.Message.Author.Role != "assistant" {
			continue
		}
		if chat_response.Message.Metadata.Timestamp == "absolute" {
			continue
		}
		tmp_fulltext := chat_response.Message.Content.Parts[0]
		chat_response.Message.Content.Parts[0] = strings.ReplaceAll(chat_response.Message.Content.Parts[0], fulltext, "")
		var delta responses.Delta = responses.Delta{
			Content: chat_response.Message.Content.Parts[0],
			Role:    "assistant",
		}
		var finish_reason interface{}
		if chat_response.Message.Metadata.FinishDetails != nil {
			finish_reason = "stop"
			delta.Content = ""
		} else {
			finish_reason = nil
		}
		completions_response := responses.ChatCompletionChunk{
			ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
			Object:  "chat.completion.chunk",
			Created: int64(chat_response.Message.CreateTime),
			Model:   "gpt-3.5-turbo-0301",
			Choices: []responses.Choices{
				{
					Index:        0,
					Delta:        delta,
					FinishReason: finish_reason,
				},
			},
		}

		// Stream the response to the client
		response_string, err := json.Marshal(completions_response)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "error parsing response",
			})
			return
		}
		if chat_request.Stream {
			_, err = c.Writer.WriteString("data: " + string(response_string) + "\n\n")
			if err != nil {
				c.JSON(500, gin.H{
					"error": "error writing response",
				})
				return
			}
		}

		// Flush the response writer buffer to ensure that the client receives each line as it's written
		c.Writer.Flush()
		fulltext = tmp_fulltext
		if finish_reason != nil {
			if !chat_request.Stream {
				full_response := responses.ChatCompletion{
					ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
					Object:  "chat.completion",
					Created: int64(chat_response.Message.CreateTime),
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
					c.JSON(500, gin.H{
						"error": "error parsing response",
					})
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
