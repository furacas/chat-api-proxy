package xyhelper

import (
	"bytes"
	"chat-api-proxy/api"
	"chat-api-proxy/common"
	"context"
	"encoding/json"
	"errors"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
	"io"
	"time"
)

type XyHelperProvider struct {
	sem *semaphore.Weighted
}

func (p *XyHelperProvider) SendRequest(c *gin.Context, originalRequest api.APIRequest) error {
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

	req, err := http.NewRequest("POST", "https://api.xyhelper.cn/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-api-xyhelper-cn-free-token-for-everyone-xyhelper")

	resp, err := common.NewClient().Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("error response code")
	}

	c.Header("X-Provider", "xyhelper")
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
