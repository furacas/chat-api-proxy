package fakeopen

import (
	"bytes"
	"chat-api-proxy/api"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
	"io"
	"net/http"
	"time"
)

type FakeOpenProvider struct {
	sem *semaphore.Weighted
}

func (p *FakeOpenProvider) SendRequest(c *gin.Context, originalRequest api.APIRequest) error {
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

	jsonData, _ := json.Marshal(originalRequest)

	req, err := http.NewRequest("POST", "https://ai.fakeopen.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer pk-this-is-a-real-free-pool-token-for-everyone")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Header(k, v)
		}
	}

	c.Header("X-Provider", "fakeopen")
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
