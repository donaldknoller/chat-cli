package common_llm

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSON(t *testing.T) {
	cm := ChatMessage{
		Role:    User,
		Content: "some text",
	}
	raw, err := json.Marshal(cm)
	assert.Nil(t, err)
	assert.Equal(t, string(raw), "{\"role\":\"user\",\"content\":\"some text\"}")

}
