package xmpp

import (
	"context"
	"encoding/xml"

	"mellium.im/xmlstream"
	"mellium.im/xmpp"
)

type Handler interface {
	HandleXMPP(
		ctx context.Context,
		t xmlstream.TokenReadEncoder,
		start *xml.StartElement,
	) (bool, error)
}

func HandleWithContext(ctx context.Context, handler Handler) xmpp.Handler {
	return &ContextualHandler{
		handler: handler,
		ctx:     ctx,
	}
}

type ContextualHandler struct {
	handler Handler
	ctx     context.Context //nolint:containedctx
}

func (c *ContextualHandler) HandleXMPP(
	t xmlstream.TokenReadEncoder,
	start *xml.StartElement,
) error {
	_, err := c.handler.HandleXMPP(c.ctx, t, start)
	return err
}
