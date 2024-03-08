package common_llm

import "encoding/json"

type StreamResponse struct {
	Err   error
	Chunk string
	Done  bool
}
type Role int

const (
	System Role = iota
	User
	Assistant
)

// roleStrings maps Role values to their string representations.
var roleStrings = []string{"system", "user", "assistant"}

// String returns the string representation of a Role.
func (r Role) String() string {
	if r < System || r > Assistant {
		return "unknown"
	}
	return roleStrings[r]
}

func (cm ChatMessage) MarshalJSON() ([]byte, error) {
	// Custom struct to define how ChatMessage should be marshaled to JSON
	type chatMessageJSON struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	// Create an instance of the custom struct with the desired fields
	msg := chatMessageJSON{
		Role:    cm.Role.String(),
		Content: cm.Content,
	}

	// Marshal the custom struct instead of the ChatMessage directly
	return json.Marshal(msg)
}

type ChatMessage struct {
	Role    `json:"role"`
	Content string `json:"content"`
}
