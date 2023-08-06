package chatgpt

import (
	"bufio"
	"bytes"
	"chat-api-proxy/api"
	"chat-api-proxy/common"
	"encoding/json"
	"errors"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
	"io"
	"strings"
)

type ChatGPTProvider struct {
	sem *semaphore.Weighted
}

func (p *ChatGPTProvider) Name() string {
	return "chatgpt"
}

func (p *ChatGPTProvider) SendRequest(c *gin.Context, originalRequest api.APIRequest) error {

	token, err := GetToken()

	if err != nil {
		return err
	}

	translatedRequest := convertAPIRequest(originalRequest)
	response, err := postConversation(translatedRequest, token)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("error sending request")
	}

	var full_response string
	for i := 3; i > 0; i-- {
		var continue_info *ContinueInfo
		var response_part string
		response_part, continue_info = Handler(c, response, token, translatedRequest, originalRequest.Stream)
		full_response += response_part
		if continue_info == nil {
			break
		}
		println("Continuing conversation")
		translatedRequest.Messages = nil
		translatedRequest.Action = "continue"
		translatedRequest.ConversationID = continue_info.ConversationID
		translatedRequest.ParentMessageID = continue_info.ParentID
		response, err = postConversation(translatedRequest, token)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "error sending request",
			})
			//return
		}
		defer response.Body.Close()
		if Handle_request_error(c, response) {
			//return
		}
	}
	if !originalRequest.Stream {
		c.JSON(200, api.NewChatCompletion(full_response))
	} else {
		c.String(200, "data: [DONE]\n\n")
	}

	return nil
}

func convertAPIRequest(apiRequest api.APIRequest) ChatGPTRequest {
	gptRequest := NewChatGPTRequest()
	if strings.HasPrefix(apiRequest.Model, "gpt-3.5") {
		gptRequest.Model = "text-davinci-002-render-sha"
	}
	for _, apiMessage := range apiRequest.Messages {
		if apiMessage.Role == "system" {
			apiMessage.Role = "critic"
		}
		gptRequest.AddMessage(apiMessage.Role, apiMessage.Content)
	}

	return gptRequest
}

func sendConversationRequest(c *gin.Context, request ChatGPTRequest, accessToken string) (*http.Response, bool) {
	jsonBytes, _ := json.Marshal(request)

	apiUrl := "https://chat.openai.com/backend-api/conversation"
	req, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := common.NewClient().Do(req)
	if err != nil {
		// c.AbortWithStatusJSON(http.StatusInternalServerError, api.ReturnMessage(err.Error()))
		return nil, true
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			// logger.Error(fmt.Sprintf(api.AccountDeactivatedErrorMessage, c.GetString(api.EmailKey)))
		}

		responseMap := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&responseMap)
		c.AbortWithStatusJSON(resp.StatusCode, responseMap)
		return nil, true
	}

	return resp, false
}

func postConversation(message ChatGPTRequest, access_token string) (*http.Response, error) {

	apiUrl := "https://chat.openai.com/backend-api/conversation"

	// JSONify the body and add it to the request
	body_json, err := json.Marshal(message)
	if err != nil {
		return &http.Response{}, err
	}

	request, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(body_json))
	if err != nil {
		return &http.Response{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	request.Header.Set("Accept", "*/*")
	if access_token != "" {
		request.Header.Set("Authorization", "Bearer "+access_token)
	}
	if err != nil {
		return &http.Response{}, err
	}
	response, err := common.NewClient().Do(request)
	return response, err
}

func Handle_request_error(c *gin.Context, response *http.Response) bool {
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var error_response map[string]interface{}
		err := json.NewDecoder(response.Body).Decode(&error_response)
		if err != nil {
			// Read response body
			body, _ := io.ReadAll(response.Body)
			c.JSON(500, gin.H{"error": gin.H{
				"message": "Unknown error",
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "500",
				"details": string(body),
			}})
			return true
		}
		c.JSON(response.StatusCode, gin.H{"error": gin.H{
			"message": error_response["detail"],
			"type":    response.Status,
			"param":   nil,
			"code":    "error",
		}})
		return true
	}
	return false
}

type ContinueInfo struct {
	ConversationID string `json:"conversation_id"`
	ParentID       string `json:"parent_id"`
}

func Handler(c *gin.Context, response *http.Response, token string, translated_request ChatGPTRequest, stream bool) (string, *ContinueInfo) {
	max_tokens := false

	// Create a bufio.Reader from the response body
	reader := bufio.NewReader(response.Body)

	// Read the response byte by byte until a newline character is encountered
	if stream {
		// Response content type is text/event-stream
		c.Header("Content-Type", "text/event-stream")
	} else {
		// Response content type is application/json
		c.Header("Content-Type", "application/json")
	}
	c.Header("X-Provider", "chatgpt")
	var finish_reason string
	var previous_text StringStruct
	var original_response ChatGPTResponse
	var isRole = true
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", nil
		}
		if len(line) < 6 {
			continue
		}
		// Remove "data: " from the beginning of the line
		line = line[6:]
		// Check if line starts with [DONE]
		if !strings.HasPrefix(line, "[DONE]") {
			// Parse the line as JSON

			err = json.Unmarshal([]byte(line), &original_response)
			if err != nil {
				continue
			}
			if original_response.Error != nil {
				c.JSON(500, gin.H{"error": original_response.Error})
				return "", nil
			}
			if original_response.Message.Author.Role != "assistant" || original_response.Message.Content.Parts == nil {
				continue
			}
			if original_response.Message.Metadata.MessageType != "next" && original_response.Message.Metadata.MessageType != "continue" || original_response.Message.EndTurn != nil {
				continue
			}
			response_string := ConvertToString(&original_response, &previous_text, isRole)
			isRole = false
			if stream {
				_, err = c.Writer.WriteString(response_string)
				if err != nil {
					return "", nil
				}
			}
			// Flush the response writer buffer to ensure that the client receives each line as it's written
			c.Writer.Flush()

			if original_response.Message.Metadata.FinishDetails != nil {
				if original_response.Message.Metadata.FinishDetails.Type == "max_tokens" {
					max_tokens = true
				}
				finish_reason = original_response.Message.Metadata.FinishDetails.Type
			}

		} else {
			if stream {
				final_line := api.StopChunk(finish_reason)
				c.Writer.WriteString("data: " + final_line.String() + "\n\n")
			}
		}
	}
	if !max_tokens {
		return previous_text.Text, nil
	}
	return previous_text.Text, &ContinueInfo{
		ConversationID: original_response.ConversationID,
		ParentID:       original_response.Message.ID,
	}
}

type StringStruct struct {
	Text string `json:"text"`
}

func ConvertToString(chatgpt_response *ChatGPTResponse, previous_text *StringStruct, role bool) string {
	translated_response := api.NewChatCompletionChunk(strings.ReplaceAll(chatgpt_response.Message.Content.Parts[0], *&previous_text.Text, ""))
	if role {
		translated_response.Choices[0].Delta.Role = chatgpt_response.Message.Author.Role
	}
	previous_text.Text = chatgpt_response.Message.Content.Parts[0]
	return "data: " + translated_response.String() + "\n\n"

}
