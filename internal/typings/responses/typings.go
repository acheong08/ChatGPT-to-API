package responses

type Message struct {
	ID         string      `json:"id"`
	Author     Author      `json:"author"`
	CreateTime float64     `json:"create_time"`
	UpdateTime interface{} `json:"update_time"`
	Content    Content     `json:"content"`
	EndTurn    interface{} `json:"end_turn"`
	Weight     float64     `json:"weight"`
	Metadata   Metadata    `json:"metadata"`
	Recipient  string      `json:"recipient"`
}

type Author struct {
	Role     string                 `json:"role"`
	Name     interface{}            `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type Metadata struct {
	Timestamp     string      `json:"timestamp_"`
	MessageType   interface{} `json:"message_type"`
	FinishDetails interface{} `json:"finish_details"`
}

type Data struct {
	Message        Message     `json:"message"`
	ConversationID string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}
type ChatCompletionChunk struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Delta        Delta       `json:"delta"`
	Index        int         `json:"index"`
	FinishReason interface{} `json:"finish_reason"`
}

type Delta struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}
