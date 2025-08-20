package ports

import (
	"context"
	"xmpp-llm-bridge/pkg/xmpp"
)

type XMPPSession interface {
	Connect(context.Context) error
	Close(context.Context) error
	Send(context.Context, xmpp.Stanza) error
	Handle(context.Context, xmpp.Handler) error
}
