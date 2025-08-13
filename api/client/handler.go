package client

import (
	"encoding/xml"
	"xmpp-llm-bridge/api/client/handlers"
	"xmpp-llm-bridge/internal/ports"

	"mellium.im/xmpp"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/stanza"
)

func NewHandler(config ports.Config, logger ports.Logger) xmpp.Handler {
	// return handlers.NewDebugHandler(logger)
	return mux.New(
		stanza.NSClient,
		mux.Message(stanza.ChatMessage, xml.Name{Local: "body"}, handlers.NewEchoMessageHandler(logger)),
		mux.Message(stanza.ChatMessage, xml.Name{Space: "http://jabber.org/protocol/chatstates"}, handlers.NewDebugHandler(logger)),
	)
}
