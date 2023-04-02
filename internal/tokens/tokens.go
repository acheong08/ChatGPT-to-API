package tokens

import (
	"sync"
)

type AccessToken struct {
	tokens []string
	lock   sync.Mutex
}

func NewAccessToken(tokens []string) AccessToken {
	return AccessToken{
		tokens: tokens,
	}
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
