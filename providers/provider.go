package providers

import (
	"chat-api-proxy/api"
	"net/http"
	"sync"
)

type Provider interface {
	SendRequest(originalRequest api.APIRequest) (*http.Response, error)
}

func PollProviders(providers []Provider, originalRequest api.APIRequest) (*http.Response, error) {
	var wg sync.WaitGroup
	responseChan := make(chan *http.Response, len(providers))
	errChan := make(chan error, len(providers))

	for _, provider := range providers {
		wg.Add(1)
		go func(p Provider) {
			defer wg.Done()
			resp, err := p.SendRequest(originalRequest)
			if err != nil {
				errChan <- err
				return
			}
			responseChan <- resp
		}(provider)
	}

	wg.Wait()
	close(responseChan)
	close(errChan)

	if len(responseChan) > 0 {
		return <-responseChan, nil
	}

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return nil, nil
}
