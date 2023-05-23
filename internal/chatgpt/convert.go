package chatgpt

import (
	typings "freechatgpt/internal/typings"
	"strings"
)

func ConvertAPIRequest(api_request typings.APIRequest) ChatGPTRequest {
	chatgpt_request := NewChatGPTRequest()
	if strings.HasPrefix(api_request.Model, "gpt-4") {
		chatgpt_request.Model = "gpt-4"
		if api_request.Model == "gpt-4-browsing" {
			chatgpt_request.Model = "gpt-4-browsing"
		}
	}
	for _, api_message := range api_request.Messages {
		if api_message.Role == "system" {
			api_message.Role = "critic"
		}
		chatgpt_request.AddMessage(api_message.Role, api_message.Content)
	}
	return chatgpt_request
}
