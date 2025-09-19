package xmpp

import (
	"context"
	"encoding/xml"
)

type StreamWriter interface {
	Write(Stanza) error
}

type Handler interface {
	HandleXMPP(ctx context.Context, t xml.TokenReader, w StreamWriter) (bool, error)
}
