package providers

import (
	"chat-api-proxy/api"
	"chat-api-proxy/providers/chatgpt"
	"chat-api-proxy/providers/fakeopen"
	"github.com/gin-gonic/gin"
)

var allProviders = []Provider{
	&fakeopen.FakeOpenProvider{},
	&chatgpt.AccountProvider{},
}

func PollProviders(c *gin.Context, originalRequest api.APIRequest) error {
	var lastError error

	for _, provider := range allProviders {
		err := provider.SendRequest(c, originalRequest)
		if err != nil {
			lastError = err
			continue
		} else {
			return nil
		}
	}

	if lastError != nil {
		return lastError
	}

	return nil
}
