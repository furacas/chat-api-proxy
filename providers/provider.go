package providers

import (
	"chat-api-proxy/api"
	"github.com/gin-gonic/gin"
)

type Provider interface {
	SendRequest(c *gin.Context, originalRequest api.APIRequest) error
}
