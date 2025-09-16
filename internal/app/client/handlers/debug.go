package handlers

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"

	"mellium.im/xmlstream"
)

// DebugHandler logs all incoming XMPP stanzas for debugging purposes
type DebugHandler struct {
	loggerProvider *providers.LoggerProvider
}

func NewDebugHandler(loggerProvider *providers.LoggerProvider) *DebugHandler {
	return &DebugHandler{
		loggerProvider: loggerProvider,
	}
}

func toString(tr xml.TokenReader) (string, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)

	for {
		tok, err := tr.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if err := enc.EncodeToken(tok); err != nil {
			return "", err
		}
	}
	if err := enc.Flush(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (h *DebugHandler) HandleXMPP(
	ctx context.Context,
	t xmlstream.TokenReadEncoder,
	start *xml.StartElement,
) (bool, error) {
	logger := h.loggerProvider.Value(ctx)
	str, err := toString(xmlstream.MultiReader(xmlstream.Token(*start), t))
	if err != nil {
		return false, err
	}

	logger.Debug("incoming", ports.Fields{
		"stanza":  start.Name.Local,
		"content": str,
	})

	return true, nil
}
