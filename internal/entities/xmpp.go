package entities

import (
	"encoding/xml"

	"mellium.im/xmpp/stanza"
)

// TODO fix entities

type MessageBody struct {
	stanza.Message
	Subject   string          `xml:"subject,omitempty"`
	Body      string          `xml:"body"`
	ChatState ChatStateActive `xml:"http://jabber.org/protocol/chatstates active,omitempty"`
}

type ComposingMessage struct {
	stanza.Message
	ChatState ChatStateComposing
}

type ChatStateActive struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates active"`
}

type ChatStateComposing struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates composing"`
}
