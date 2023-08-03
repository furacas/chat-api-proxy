package main

import (
	"chat-api-proxy/api"
	"chat-api-proxy/providers"
	"github.com/gin-gonic/gin"
	"io"
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

	allProviders := []providers.Provider{
		&providers.FakeOpenProvider{},
		// add other allProviders here
	}

	resp, err := providers.PollProviders(allProviders, originalRequest)
	if err != nil || resp == nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Header(k, v)
		}
	}

	c.Status(resp.StatusCode)

	buf := make([]byte, 256) // 1 byte buffer
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, err := c.Writer.Write(buf[:n])
			if err != nil {
				// Handle error.
				return
			}
			c.Writer.Flush()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			// Handle error.
			return
		}
	}

}
