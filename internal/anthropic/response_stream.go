package anthropic

import "github.com/donaldknoller/chat-cli/internal/common_llm"

type StreamResponse struct {
	Err   error
	Chunk string
	Done  bool
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

func (r *StreamResponse) ToCommon() common_llm.StreamResponse {
	return common_llm.StreamResponse{
		Chunk: r.Chunk,
		Err:   r.Err,
		Done:  r.Done,
	}
}
