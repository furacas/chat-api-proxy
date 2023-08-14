package fakeopen

import (
	"chat-api-proxy/api"
	"chat-api-proxy/common"
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
	"time"
)

type FakeOpenProvider struct {
	sem *semaphore.Weighted
}

func (p *FakeOpenProvider) Name() string {
	return "fakeopen"
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

	return common.SendRequest(c, originalRequest, "https://ai.fakeopen.com/v1/chat/completions", "pk-this-is-a-real-free-pool-token-for-everyone", p.Name())
}
