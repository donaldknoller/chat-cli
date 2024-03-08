package common_llm

type Session interface {
	AddMessage(string, Role) Session
	AddChunk(string)
	PostDataStream(chan StreamResponse)
	IsStreaming() bool
	StreamingMessage() string
	ChatMessages() []ChatMessage
}
