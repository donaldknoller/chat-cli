package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/donaldknoller/chat-cli/internal/config"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type LLMClient struct {
	HttpClient *http.Client
}

func NewClient() *LLMClient {
	return &LLMClient{
		HttpClient: &http.Client{},
	}
}

func (c *LLMClient) PostData(request Request) (string, error) {

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// Create the POST request
	req, err := http.NewRequest("POST", viper.GetString(config.LLM_API_HOST), bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", viper.GetString(config.LLM_API_KEY))
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if viper.GetBool(config.DEBUG) {
		fmt.Println("messages: ", request.Messages)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response Response
	if responseErr := json.Unmarshal(body, &response); responseErr != nil {
		return "", responseErr
	}
	if viper.GetBool(config.DEBUG) {
		fmt.Println("content:", response.Content)
	}
	if len(response.Content) < 1 {
		return "", fmt.Errorf("len(response) has value: %d", len(response.Content))
	}

	return response.Content[0].Text, nil
}
