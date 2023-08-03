package chatgpt

import (
	"errors"
	"sync"
)

var tokenManager = initializeTokens([]string{})

type TokenManager struct {
	tokens []string
	lock   sync.Mutex
}

func initializeTokens(tokens []string) *TokenManager {
	return &TokenManager{
		tokens: tokens,
	}
}

func GetToken() (string, error) {
	tokenManager.lock.Lock()
	defer tokenManager.lock.Unlock()

	if len(tokenManager.tokens) == 0 {
		return "", errors.New("no token find")
	}

	token := tokenManager.tokens[0]
	tokenManager.tokens = append(tokenManager.tokens[1:], token)
	return token, nil
}
