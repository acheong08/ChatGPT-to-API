package chatgpt

import (
	chatgpt_types "freechatgpt/typings/chatgpt"
	official_types "freechatgpt/typings/official"
	"strings"
)

func ConvertAPIRequest(api_request official_types.APIRequest) chatgpt_types.ChatGPTRequest {
	chatgpt_request := chatgpt_types.NewChatGPTRequest()
	if strings.HasPrefix(api_request.Model, "gpt-4") {
		chatgpt_request.Model = "gpt-4"
		if api_request.Model == "gpt-4-browsing" || api_request.Model == "gpt-4-plugins" || api_request.Model == "gpt-4-mobile" || api_request.Model == "gpt-4-code-interpreter" {
			chatgpt_request.Model = api_request.Model
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
