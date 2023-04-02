package types

type APIRequest struct {
	Messages []api_message `json:"messages"`
}

type api_message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
