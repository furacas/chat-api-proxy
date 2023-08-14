package chatanywhere

import (
	"chat-api-proxy/api"
	"chat-api-proxy/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/sync/semaphore"
	"os"
	"time"
)

type ChatAnyWhereProvider struct {
	sem *semaphore.Weighted
}

func (p *ChatAnyWhereProvider) Name() string {
	return "chatanywhere"
}

func (p *ChatAnyWhereProvider) SendRequest(c *gin.Context, originalRequest api.APIRequest) error {
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

	return common.SendRequest(c, originalRequest, "https://api.chatanywhere.cn/v1/chat/completions", os.Getenv("CHATANYWHERE_KEY"), p.Name())
}
