package main

import (
	"chat-api-proxy/api"
	"chat-api-proxy/providers"
	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func optionsHandler(c *gin.Context) {
	// Set headers for CORS
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func completionsHandler(c *gin.Context) {

	var originalRequest api.APIRequest
	err := c.BindJSON(&originalRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": gin.H{
			"message": "Request must be proper JSON",
			"type":    "invalid_request_error",
			"param":   nil,
			"code":    err.Error(),
		}})
		return
	}

	err = providers.PollProviders(c, originalRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "no providers available"})
		return
	}

}

func providerStatHandler(c *gin.Context) {
	snapshot := make(map[string]*providers.ProviderStat)
	providers.ProviderStats.Range(func(key, value interface{}) bool {
		snapshot[key.(string)] = value.(*providers.ProviderStat)
		return true
	})

	c.JSON(200, snapshot)
}
