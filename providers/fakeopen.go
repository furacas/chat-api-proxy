package providers

import (
	"bytes"
	"chat-api-proxy/api"
	"encoding/json"
	"net/http"
)

type FakeOpenProvider struct {
}

func (f *FakeOpenProvider) SendRequest(originalRequest api.APIRequest) (*http.Response, error) {
	jsonData, _ := json.Marshal(originalRequest)

	req, err := http.NewRequest("POST", "https://ai.fakeopen.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer pk-this-is-a-real-free-pool-token-for-everyone")

	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
