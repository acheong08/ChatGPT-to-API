package types

import (
	"math/rand"
)

type ChatMessage struct {
	Action    string   `json:"action"`
	Model     string   `json:"model"`
	Jailbreak string   `json:"jailbreak"`
	Meta      metadata `json:"meta"`
}

type metadata struct {
	ID      int64           `json:"id"`
	Content message_content `json:"content"`
}

type message_content struct {
	Conversation   []conversation_message `json:"conversation"`
	InternetAccess bool                   `json:"internet_access"`
	ContentType    string                 `json:"content_type"`
	Parts          []message_part         `json:"parts"`
}

type conversation_message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type message_part struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func generateRandomInt64() int64 {
	return rand.Int63()
}

func NewChatMessage() ChatMessage {
	return ChatMessage{
		Action:    "_ask",
		Model:     "text-gpt-0004-render-sha-0314",
		Jailbreak: "default",
		Meta: metadata{
			ID: generateRandomInt64(),
			Content: message_content{
				Conversation:   []conversation_message{},
				InternetAccess: false,
				ContentType:    "text",
				Parts:          []message_part{},
			},
		},
	}
}

func (chat *ChatMessage) AddMessage(role string, content string) {
	chat.Meta.Content.Conversation = append(chat.Meta.Content.Conversation, conversation_message{
		Role:    role,
		Content: content,
	})
}

func (chat *ChatMessage) AddPart(role string, content string) {
	chat.Meta.Content.Parts = append(chat.Meta.Content.Parts, message_part{
		Role:    role,
		Content: content,
	})
}

func (chat *ChatMessage) Rollback() {
	chat.Meta.Content.Conversation = chat.Meta.Content.Conversation[:len(chat.Meta.Content.Conversation)-1]
}
