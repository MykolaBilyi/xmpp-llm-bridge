package xml

import (
	"encoding/xml"
	"io"
)

func DuplicateReader(t xml.TokenReader) (xml.TokenReader, xml.TokenReader, error) {
	tokens := make([]xml.Token, 0)
	copy := make([]xml.Token, 0)
	for {
		token, err := t.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		tokens = append(tokens, xml.CopyToken(token))
		copy = append(copy, xml.CopyToken(token))
	}

	return &tokenReaderFromSlice{tokens: tokens, index: 0}, &tokenReaderFromSlice{tokens: copy, index: 0}, nil
}

// tokenReaderFromSlice implements xml.TokenReader from a slice of tokens
type tokenReaderFromSlice struct {
	tokens []xml.Token
	index  int
}

func (tr *tokenReaderFromSlice) Token() (xml.Token, error) {
	if tr.index >= len(tr.tokens) {
		return nil, io.EOF
	}
	token := tr.tokens[tr.index]
	tr.index++
	return token, nil
}
