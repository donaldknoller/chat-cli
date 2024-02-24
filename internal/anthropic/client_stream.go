package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/donaldknoller/chat-cli/internal/config"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func (c *LLMClient) PostDataStream(request Request, responseChan chan StreamResponse) {
	defer close(responseChan)
	requestBody, err := json.Marshal(request)
	if err != nil {
		responseChan <- StreamResponse{Err: err}
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		viper.GetString(config.LLM_API_HOST),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		responseChan <- StreamResponse{Err: err}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", viper.GetString(config.LLM_API_KEY))
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		responseChan <- StreamResponse{Err: err}
		return
	}
	defer resp.Body.Close()
	if viper.GetBool(config.DEBUG) {
		fmt.Println("messages: ", request.Messages)
	}
	scanner := bufio.NewScanner(resp.Body)
	eventType := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "event:") {
			eventType = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "data:") {
			dataJSON := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if viper.GetBool(config.DEBUG) {
				fmt.Println("dataChunk: ", dataJSON)
			}
			response := parseStreamingDataChunk(dataJSON, eventType)
			responseChan <- response
		}
	}

}

func parseStreamingDataChunk(rawJSON, eventType string) StreamResponse {
	switch eventType {
	// do nothing for irrelevant messages, but note them in case
	case "message_start", "message_delta", "message_stop", "content_block_stop", "ping":
		return StreamResponse{}
	case "content_block_start":
		var block ContentBlockEvent
		err := json.Unmarshal([]byte(rawJSON), &block)
		if err != nil {
			return StreamResponse{
				Err: err,
			}
		}
		return StreamResponse{
			Chunk: block.ContentBlock.Text,
		}
	case "content_block_delta":
		var block ContentBlockDeltaEvent
		err := json.Unmarshal([]byte(rawJSON), &block)
		if err != nil {
			return StreamResponse{
				Err: err,
			}
		}
		return StreamResponse{
			Chunk: block.Delta.Text,
		}
	default:
		return StreamResponse{}
	}
}
