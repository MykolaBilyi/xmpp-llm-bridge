package xml_test

import (
	"encoding/xml"
	"strings"
	"testing"

	myxml "xmpp-llm-bridge/pkg/xml"
)

func TestExtractStartElement(t *testing.T) {
	tests := []struct {
		name    string
		xmlData string
		want    *xml.StartElement
		wantErr bool
	}{
		{
			name:    "simple element",
			xmlData: "<message>content</message>",
			want: &xml.StartElement{
				Name: xml.Name{Local: "message"},
			},
			wantErr: false,
		},
		{
			name:    "element with attributes",
			xmlData: `<message type="chat" from="user@example.com">content</message>`,
			want: &xml.StartElement{
				Name: xml.Name{Local: "message"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "type"}, Value: "chat"},
					{Name: xml.Name{Local: "from"}, Value: "user@example.com"},
				},
			},
			wantErr: false,
		},
		{
			name:    "element with namespace",
			xmlData: `<message xmlns="jabber:client">content</message>`,
			want: &xml.StartElement{
				Name: xml.Name{Space: "jabber:client", Local: "message"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "xmlns"}, Value: "jabber:client"},
				},
			},
			wantErr: false,
		},
		{
			name:    "self-closing element",
			xmlData: `<presence type="unavailable"/>`,
			want: &xml.StartElement{
				Name: xml.Name{Local: "presence"},
				Attr: []xml.Attr{
					{Name: xml.Name{Local: "type"}, Value: "unavailable"},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty string",
			xmlData: "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "invalid xml",
			xmlData: "<invalid",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "text only",
			xmlData: "just text",
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO test copy of reader
			reader := strings.NewReader(tt.xmlData)
			_, got, err := myxml.ExtractStartElement(xml.NewDecoder(reader))

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractStartElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.want == nil {
					if got != nil {
						t.Errorf("ExtractStartElement() got = %v, want nil", got)
					}
					return
				}

				if got.Name != tt.want.Name {
					t.Errorf("ExtractStartElement() name = %v, want %v", got.Name, tt.want.Name)
				}

				if len(got.Attr) != len(tt.want.Attr) {
					t.Errorf("ExtractStartElement() attrs length = %d, want %d", len(got.Attr), len(tt.want.Attr))
				}

				for i, attr := range got.Attr {
					if i < len(tt.want.Attr) && attr != tt.want.Attr[i] {
						t.Errorf("ExtractStartElement() attr[%d] = %v, want %v", i, attr, tt.want.Attr[i])
					}
				}
			}
		})
	}
}
