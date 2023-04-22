package chatgpt

import (
	typings "freechatgpt/internal/typings"
)

func ConvertAPIRequest(api_request typings.APIRequest) typings.ChatMessage {
	chatgpt_request := typings.NewChatMessage()
	for _, api_message := range api_request.Messages {
		chatgpt_request.AddMessage(api_message.Role, api_message.Content)
	}
	return chatgpt_request
}
