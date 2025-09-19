package stanza

type Stanza string

const Message Stanza = "message"

type MessageType string

const (
	NormalMessage MessageType = "normal"
	ChatMessage   MessageType = "chat"
)
