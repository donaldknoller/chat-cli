package anthropic

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Messages    []Message `json:"messages"`
	Stream      *bool     `json:"stream,omitempty"`
	Temperature *float64  `json:"temperature,omitempty"`
	TopP        *float64  `json:"top_p,omitempty"`
	TopK        *float64  `json:"top_k,omitempty"`
}

func (r *Request) Merge(responseMessage string) *Request {
	r.Messages = append(r.Messages, Message{
		Role:    "assistant",
		Content: responseMessage,
	})
	return r
}
func (r *Request) InsertMessage(requestMessage string) *Request {
	r.Messages = append(r.Messages, Message{
		Role:    "user",
		Content: requestMessage,
	})
	return r
}
