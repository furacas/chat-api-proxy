package xyhelper

import (
	"chat-api-proxy/api"
	"chat-api-proxy/common"
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
	"time"
)

type XyHelperProvider struct {
	sem *semaphore.Weighted
}

func (p *XyHelperProvider) Name() string {
	return "xyhelper"
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

	return common.SendRequest(c, originalRequest, "https://api.xyhelper.cn/v1/chat/completions", "sk-api-xyhelper-cn-free-token-for-everyone-xyhelper")
}
