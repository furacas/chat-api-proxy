package providers

import (
	"chat-api-proxy/api"
	"chat-api-proxy/providers/ava"
	"chat-api-proxy/providers/chatanywhere"
	"chat-api-proxy/providers/chatgpt"
	"chat-api-proxy/providers/fakeopen"
	"chat-api-proxy/providers/xyhelper"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"os"
	"time"
)

type Provider interface {
	SendRequest(c *gin.Context, originalRequest api.APIRequest) error

	Name() string
}

var allProviders []Provider

func init() {
	if os.Getenv("FAKEOPEN_ENABLED") != "false" {
		log.Printf("FakeOpenProvider enabled")
		allProviders = append(allProviders, &fakeopen.FakeOpenProvider{})
	}
	if os.Getenv("XYHELPER_ENABLED") != "false" {
		log.Printf("XyHelperProvider enabled")
		allProviders = append(allProviders, &xyhelper.XyHelperProvider{})
	}
	if os.Getenv("CHATGPT_ENABLED") != "false" {
		log.Printf("ChatGPTProvider enabled")
		allProviders = append(allProviders, &chatgpt.ChatGPTProvider{})
	}
	if os.Getenv("AVA_ENABLED") != "false" {
		log.Printf("AvaProvider enabled")
		allProviders = append(allProviders, &ava.AvaProvider{})
	}

	if os.Getenv("CHATANYWHERE_KEY") != "" {
		log.Printf("ChatAnyWhereProvider enabled")
		allProviders = append(allProviders, &chatanywhere.ChatAnyWhereProvider{})
	}
}

func PollProviders(c *gin.Context, originalRequest api.APIRequest) error {
	specifiedProvider := c.GetHeader("X-Provider")

	if specifiedProvider != "" {
		// Iterate over all providers to find a match
		for _, provider := range allProviders {
			if provider.Name() == specifiedProvider {
				err := provider.SendRequest(c, originalRequest)
				if err != nil {
					return err
				}
				updateStat(specifiedProvider, true)
				return nil
			}
		}
		// If specified provider isn't found among all providers
		return fmt.Errorf("Specified provider '%s' not found", specifiedProvider)
	}

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	indices := r.Perm(len(allProviders))

	var lastError error

	for _, index := range indices {
		provider := allProviders[index]

		providerName := provider.Name()

		err := provider.SendRequest(c, originalRequest)
		if err != nil {
			lastError = err
			updateStat(providerName, false)
			continue
		} else {
			updateStat(providerName, true)
			return nil
		}
	}

	if lastError != nil {
		return lastError
	}

	return nil
}
