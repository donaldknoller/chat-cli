package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/donaldknoller/chat-cli/internal/common_llm"
	"github.com/donaldknoller/chat-cli/internal/config"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

type AnthropicSession struct {
	*Request
	*LLMClient
	isStreaming      bool
	streamingMessage string
}

func NewSession() common_llm.Session {
	useStream := true
	request := &Request{
		Model:     viper.GetString(config.LLM_API_MODEL),
		MaxTokens: 1024,
		Messages:  []common_llm.ChatMessage{},
		Stream:    &useStream,
	}
	client := NewClient()
	session := AnthropicSession{
		Request:   request,
		LLMClient: client,
	}
	return &session
}

func (s *AnthropicSession) AddMessage(message string, role common_llm.Role) common_llm.Session {
	s.Request = s.Request.AddMessage(message, role)
	return s
}
func (s *AnthropicSession) AddChunk(chunk string) {
	if s.isStreaming {
		s.streamingMessage = s.streamingMessage + chunk
	}
}
func (s *AnthropicSession) PostDataStream(responseChan chan common_llm.StreamResponse) {
	s.isStreaming = true
	s.streamingMessage = ""
	defer close(responseChan)
	requestBody, err := json.Marshal(s.Request)
	if err != nil {
		responseChan <- common_llm.StreamResponse{Err: err}
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		viper.GetString(config.LLM_API_HOST),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		responseChan <- common_llm.StreamResponse{Err: err}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", viper.GetString(config.LLM_API_KEY))
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.LLMClient.HttpClient.Do(req)
	if err != nil {
		responseChan <- common_llm.StreamResponse{Err: err}
		return
	}
	defer resp.Body.Close()
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
			responseChan <- response.ToCommon()
		}
	}
	s.isStreaming = false
	s.AddMessage(s.streamingMessage, common_llm.Assistant)
	s.streamingMessage = ""
	responseChan <- common_llm.StreamResponse{
		Done: true,
	}
}
func (s *AnthropicSession) IsStreaming() bool {
	return s.isStreaming
}
func (s *AnthropicSession) StreamingMessage() string {
	return s.streamingMessage
}
func (s *AnthropicSession) ChatMessages() []common_llm.ChatMessage {
	return s.Messages
}
