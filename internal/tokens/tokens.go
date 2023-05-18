package tokens

import (
	"encoding/json"
	"os"
	"sync"
)

type AccessToken struct {
	tokens []string
	lock   sync.Mutex
}

func NewAccessToken(tokens []string, save bool) AccessToken {
	// Save the tokens to a file
	if _, err := os.Stat("access_tokens.json"); os.IsNotExist(err) {
		// Create the file
		file, err := os.Create("access_tokens.json")
		if err != nil {
			return AccessToken{}
		}
		defer file.Close()
	}
	if save {
		saved := Save(tokens)
		if saved == false {
			return AccessToken{}
		}
	}
	return AccessToken{
		tokens: tokens,
	}
}

func Save(tokens []string) bool {
	file, err := os.OpenFile("access_tokens.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return false
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(tokens)
	if err != nil {
		return false
	}
	return true
}

func (a *AccessToken) GetToken() string {
	a.lock.Lock()
	defer a.lock.Unlock()

	if len(a.tokens) == 0 {
		return ""
	}

	token := a.tokens[0]
	a.tokens = append(a.tokens[1:], token)
	return token
}
