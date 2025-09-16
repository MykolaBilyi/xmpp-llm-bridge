package mux

import (
	"context"
	"encoding/xml"
	"xmpp-llm-bridge/pkg/xmpp"

	"mellium.im/xmlstream"
)

type mux struct {
	ns       string
	handlers []xmpp.Handler
}

func New(ns string, handlers ...xmpp.Handler) xmpp.Handler {
	return &mux{
		ns:       ns,
		handlers: handlers,
	}
}

func (m *mux) HandleXMPP(
	ctx context.Context,
	t xmlstream.TokenReadEncoder,
	start *xml.StartElement,
) (bool, error) {
	if start.Name.Space != m.ns {
		return false, nil // Not our namespace, skip
	}

	for _, handler := range m.handlers {
		h, err := handler.HandleXMPP(ctx, t, start)
		if err != nil {
			return false, err
		}
		if h {
			return true, nil
		}
	}
	return false, nil
}
