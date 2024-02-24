package anthropic

type Content struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type Response struct {
	Content      []Content `json:"content"`
	ID           string    `json:"id"`
	Model        string    `json:"model"`
	Role         string    `json:"role"`
	StopReason   string    `json:"stop_reason"`
	StopSequence *string   `json:"stop_sequence"` // Using a pointer to represent nullability
	Type         string    `json:"type"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}
