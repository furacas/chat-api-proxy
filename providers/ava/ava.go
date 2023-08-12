package ava

import (
	"bufio"
	"bytes"
	"chat-api-proxy/api"
	"chat-api-proxy/common"
	"encoding/json"
	"errors"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/sync/semaphore"
	"io"
	"strings"
	"time"
)

type AvaProvider struct {
	sem *semaphore.Weighted
}

func (p *AvaProvider) Name() string {
	return "ava"
}

type AvaRequest struct {
	Stream      bool             `json:"stream"`
	Temperature float32          `json:"temperature"`
	Messages    []api.APIMessage `json:"messages"`
}

func (p *AvaProvider) SendRequest(c *gin.Context, originalRequest api.APIRequest) error {
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

	avaRequest := convertRequest(originalRequest)

	jsonData, _ := json.Marshal(avaRequest)

	req, err := http.NewRequest("POST", "https://ava-alpha-api.codelink.io/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := common.NewClient().Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("error response code")
	}

	c.Header("X-Provider", "ava")
	c.Status(resp.StatusCode)

	if originalRequest.Stream {
		c.Header("Content-Type", "text/event-stream")

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
	} else {
		c.Header("Content-Type", "application/json")

		var messageBuilder strings.Builder

		finalResponse := map[string]interface{}{}
		initial := true // 用于判断是否为第一条消息，以捕获所有的常量字段

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				// Handle error.
				return err
			}

			if len(line) > 0 && strings.HasPrefix(line, "data: ") {
				jsonStr := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
				if jsonStr != "[DONE]" {
					var data map[string]interface{}
					if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
						// Handle JSON parsing error.
						return err
					}

					// 如果是第一条消息，复制所有字段到finalResponse
					if initial {
						for key, value := range data {
							finalResponse[key] = value
						}
						initial = false
					}

					// 提取delta.content并拼接到消息中
					if choices, ok := data["choices"].([]interface{}); ok {
						for _, choice := range choices {
							choiceMap, isMap := choice.(map[string]interface{})
							if !isMap {
								continue
							}

							if delta, ok := choiceMap["delta"].(map[string]interface{}); ok {
								if content, ok := delta["content"].(string); ok {
									messageBuilder.WriteString(content)
								}
							}
						}
					}
				} else {
					// 当读到 "[DONE]" 时，停止读取
					break
				}
			}
		}
		// 将拼接的消息添加到最终的响应中
		if choices, ok := finalResponse["choices"].([]interface{}); ok && len(choices) > 0 {
			choiceMap, isMap := choices[0].(map[string]interface{})
			if isMap {
				if delta, ok := choiceMap["delta"].(map[string]interface{}); ok {
					delta["content"] = messageBuilder.String()
				}
			}
		}

		// 把最终的响应转换为JSON并写入到响应
		jsonData, err := json.Marshal(finalResponse)
		if err != nil {
			// Handle marshalling error.
			return err
		}
		_, err = c.Writer.Write(jsonData)
		if err != nil {
			// Handle error.
			return err
		}

	}

	return nil

}

func convertRequest(request api.APIRequest) AvaRequest {
	ava := AvaRequest{
		Stream:      true,
		Temperature: 0.6,
	}

	ava.Messages = request.Messages

	return ava
}
