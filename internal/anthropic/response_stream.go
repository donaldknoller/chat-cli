package anthropic

type StreamResponse struct {
	Err   error
	Chunk string
}

type ContentBlockEvent struct {
	Type         string       `json:"type"`
	Index        int          `json:"index"`
	ContentBlock ContentBlock `json:"content_block"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ContentBlockDeltaEvent struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta Delta  `json:"delta"`
}

type Delta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
