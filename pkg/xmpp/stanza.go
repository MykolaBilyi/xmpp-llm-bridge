package xmpp

type Stanza struct {
	message  interface{}
	presence interface{}
	iq       interface{}
}

func (s Stanza) IsMessage() bool {
	return s.message != nil
}

func (s Stanza) Get() interface{} {
	if s.IsMessage() {
		return s.message
	}
	return nil
}

// TODO Maybe this is overkill
func Message(msg any) Stanza {
	return Stanza{
		message: msg,
	}
}
