package providers

import (
	"chat-api-proxy/api"
	"chat-api-proxy/providers/chatgpt"
	"chat-api-proxy/providers/fakeopen"
	"chat-api-proxy/providers/xyhelper"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

var allProviders = []Provider{
	&fakeopen.FakeOpenProvider{},
	&xyhelper.XyHelperProvider{},
	&chatgpt.AccountProvider{},
}

func PollProviders(c *gin.Context, originalRequest api.APIRequest) error {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	indices := r.Perm(len(allProviders))

	var lastError error

	for _, index := range indices {
		provider := allProviders[index]

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
