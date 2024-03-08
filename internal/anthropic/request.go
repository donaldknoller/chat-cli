package anthropic

import "github.com/donaldknoller/chat-cli/internal/common_llm"

type Request struct {
	Model       string                   `json:"model"`
	MaxTokens   int                      `json:"max_tokens"`
	Messages    []common_llm.ChatMessage `json:"messages"`
	Stream      *bool                    `json:"stream,omitempty"`
	Temperature *float64                 `json:"temperature,omitempty"`
	TopP        *float64                 `json:"top_p,omitempty"`
	TopK        *float64                 `json:"top_k,omitempty"`
}

func (r *Request) AddMessage(responseMessage string, role common_llm.Role) *Request {
	r.Messages = append(r.Messages, common_llm.ChatMessage{
		Role:    role,
		Content: responseMessage,
	})
	return r
}
