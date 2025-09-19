package xml

import (
	"encoding/xml"
	"io"
)

func ExtractStartElement(reader xml.TokenReader) (xml.TokenReader, *xml.StartElement, error) {
	var tokens []xml.Token
	var startElement *xml.StartElement

	for {
		token, err := reader.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		tokens = append(tokens, xml.CopyToken(token))

		if se, ok := token.(xml.StartElement); ok && startElement == nil {
			startElement = &se
		}
	}

	return &tokenReaderFromSlice{tokens: tokens, index: 0}, startElement, nil
}
