package mux

import (
	"context"
	"encoding/xml"
	myxml "xmpp-llm-bridge/pkg/xml"
	"xmpp-llm-bridge/pkg/xmpp"
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

func (m *mux) HandleXMPP(ctx context.Context, t xml.TokenReader, w xmpp.StreamWriter) (bool, error) {
	t, start, err := myxml.ExtractStartElement(t)
	if err != nil {
		return false, err
	}
	if start.Name.Space != m.ns {
		return false, nil // Not our namespace, skip
	}

	for _, handler := range m.handlers {
		tokens, copy, err := myxml.DuplicateReader(t)
		if err != nil {
			return false, err
		}
		t = tokens
		h, err := handler.HandleXMPP(ctx, copy, w)
		if err != nil {
			return false, err
		}
		if h {
			return true, nil
		}
	}
	return false, nil
}
