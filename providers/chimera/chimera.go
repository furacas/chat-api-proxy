package chimera

import (
	"chat-api-proxy/api"
	"chat-api-proxy/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/sync/semaphore"
	"os"
	"time"
)

type ChimeraProvider struct {
	sem *semaphore.Weighted
}

func (p *ChimeraProvider) Name() string {
	return "chimera"
}

func (p *ChimeraProvider) SendRequest(c *gin.Context, originalRequest api.APIRequest) error {
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

	return common.SendRequest(c, originalRequest, "https://chimeragpt.adventblocks.cc/api/v1/chat/completions", os.Getenv("CHIMERA_KEY"), p.Name())
}
