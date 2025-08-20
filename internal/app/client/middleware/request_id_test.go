package middleware

import (
	"context"
	"encoding/xml"
	"testing"
	"xmpp-llm-bridge/internal/entities"
	"xmpp-llm-bridge/internal/providers"

	"mellium.im/xmlstream"
)

// MockHandler implements xmpp.Handler for testing
type MockHandler struct {
	handleFunc func(ctx context.Context, t xmlstream.TokenReadEncoder, start *xml.StartElement) (bool, error)
}

func (m *MockHandler) HandleXMPP(ctx context.Context, t xmlstream.TokenReadEncoder, start *xml.StartElement) (bool, error) {
	if m.handleFunc != nil {
		return m.handleFunc(ctx, t, start)
	}
	return true, nil
}

func TestRequestIDMiddleware(t *testing.T) {
	requestIdProvider := providers.NewRequestIdProvider()

	// Track if the handler was called with a request ID
	var capturedRequestID *entities.RequestId
	mockHandler := &MockHandler{
		handleFunc: func(ctx context.Context, t xmlstream.TokenReadEncoder, start *xml.StartElement) (bool, error) {
			capturedRequestID = requestIdProvider.Value(ctx)
			return true, nil
		},
	}

	// FIXME
	middleware := WithRequestID(mockHandler, requestIdProvider, nil)

	// Create a test XML element
	start := &xml.StartElement{
		Name: xml.Name{Local: "test"},
	}

	// Call the middleware
	ctx := context.Background()
	handled, err := middleware.HandleXMPP(ctx, nil, start)

	// Verify the results
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !handled {
		t.Error("Expected handler to return true")
	}

	if capturedRequestID == nil {
		t.Error("Expected request ID to be set in context")
	} else if *capturedRequestID == "" {
		t.Error("Expected request ID to be non-empty")
	}
}

func TestRequestIDGeneration(t *testing.T) {
	// Generate multiple request IDs and ensure they're unique
	ids := make(map[entities.RequestId]bool)

	for i := 0; i < 100; i++ {
		id := generateRequestID()
		if ids[id] {
			t.Errorf("Generated duplicate request ID: %s", id)
		}
		ids[id] = true

		// Verify the ID is not empty and has reasonable length
		if len(string(id)) == 0 {
			t.Error("Generated empty request ID")
		}
		if len(string(id)) != 16 { // 8 bytes = 16 hex characters
			t.Errorf("Expected request ID length of 16, got %d", len(string(id)))
		}
	}
}
