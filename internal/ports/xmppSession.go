package ports

import (
	"context"
	"encoding/xml"
	"xmpp-llm-bridge/pkg/xmpp"
)

type XMPPSession interface {
	Connect(context.Context) error
	Close(context.Context) error
	// TODO Think of a better argument type
	Send(context.Context, xml.TokenReader) error
	Handle(context.Context, xmpp.Handler) error
}
