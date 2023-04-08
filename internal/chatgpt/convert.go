package chatgpt

import (
	typings "freechatgpt/internal/typings"
)

func ConvertAPIRequest(api_request typings.APIRequest) typings.ChatGPTRequest {
	chatgpt_request := typings.NewChatGPTRequest()
	for _, api_message := range api_request.Messages {
		if api_message.Role == "system" {
			api_message.Role = "user"
			api_message.Content = "system: " + api_message.Content
		}
		chatgpt_request.AddMessage(api_message.Role, api_message.Content)
	}
	return chatgpt_request
}
