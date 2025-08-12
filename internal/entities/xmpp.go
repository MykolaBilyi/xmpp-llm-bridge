package entities

import "mellium.im/xmpp/stanza"

type MessageBody struct {
	stanza.Message
	Subject string `xml:"subject,omitempty"`
	Body    string `xml:"body"`
}
