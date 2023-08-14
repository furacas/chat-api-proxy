package common

import (
	"bytes"
	"chat-api-proxy/api"
	"encoding/json"
	"errors"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"io"
)

func SendRequest(c *gin.Context, originalRequest api.APIRequest, apiEnterPoint string, sk string) error {

	jsonData, _ := json.Marshal(originalRequest)

	req, err := http.NewRequest("POST", apiEnterPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sk)

	resp, err := NewClient().Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("error response code")
	}

	if originalRequest.Stream {
		c.Header("Content-Type", "text/event-stream")
	} else {
		c.Header("Content-Type", "application/json")
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
