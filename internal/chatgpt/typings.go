package chatgpt

import "github.com/google/uuid"

type chatgpt_message struct {
	ID      uuid.UUID       `json:"id"`
	Author  chatgpt_author  `json:"author"`
	Content chatgpt_content `json:"content"`
}

type chatgpt_content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type chatgpt_author struct {
	Role string `json:"role"`
}

type ChatGPTRequest struct {
	Action          string            `json:"action"`
	Messages        []chatgpt_message `json:"messages"`
	ParentMessageID string            `json:"parent_message_id,omitempty"`
	Model           string            `json:"model"`
}

func NewChatGPTRequest() ChatGPTRequest {
	return ChatGPTRequest{
		Action:          "next",
		ParentMessageID: uuid.NewString(),
		Model:           "text-davinci-002-render-sha",
	}
}

func (c *ChatGPTRequest) AddMessage(role string, content string) {
	c.Messages = append(c.Messages, chatgpt_message{
		ID:      uuid.New(),
		Author:  chatgpt_author{Role: role},
		Content: chatgpt_content{ContentType: "text", Parts: []string{content}},
	})
}
