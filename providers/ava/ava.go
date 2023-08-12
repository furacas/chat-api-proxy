package ava

import (
	"bytes"
	"chat-api-proxy/api"
	"chat-api-proxy/common"
	"encoding/json"
	"errors"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/sync/semaphore"
	"io"
	"time"
)

type AvaProvider struct {
	sem *semaphore.Weighted
}

func (p *AvaProvider) Name() string {
	return "ava"
}

type AvaRequest struct {
	Stream      bool             `json:"stream"`
	Temperature float32          `json:"temperature"`
	Messages    []api.APIMessage `json:"messages"`
}

func (p *AvaProvider) SendRequest(c *gin.Context, originalRequest api.APIRequest) error {
	if p.sem == nil {
		p.sem = semaphore.NewWeighted(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := p.sem.Acquire(ctx, 1)

	if err != nil {
		return err
	}
	defer p.sem.Release(1)

	avaRequest := convertRequest(originalRequest)

	jsonData, _ := json.Marshal(avaRequest)

	req, err := http.NewRequest("POST", "https://ava-alpha-api.codelink.io/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := common.NewClient().Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("error response code")
	}

	if originalRequest.Stream {
		c.Header("Content-Type", "text/event-stream")
	} else {
		c.Header("Content-Type", "application/json")
	}

	c.Header("X-Provider", "ava")
	c.Status(resp.StatusCode)

	buf := make([]byte, 256) // 1 byte buffer
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err := c.Writer.Write(buf[:n])
			if err != nil {
				// Handle error.
				return err
			}
			c.Writer.Flush()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			// Handle error.
			return err
		}
	}

	return nil

}

func convertRequest(request api.APIRequest) AvaRequest {
	ava := AvaRequest{
		Stream:      request.Stream,
		Temperature: 0.6,
	}

	ava.Messages = request.Messages

	return ava
}
