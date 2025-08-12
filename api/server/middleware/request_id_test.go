package middleware_test

import (
	"__template__/api/server/middleware"
	"__template__/internal/adapters"
	"__template__/internal/entities"
	"__template__/internal/ports"
	"__template__/internal/providers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type requestIdMiddlewareTestSetup struct {
	config            ports.Config
	requestIdProvider providers.RequestIdProvider
	responseRecorder  *httptest.ResponseRecorder
}

func setupRequestIdMiddlewareTest(t *testing.T) *requestIdMiddlewareTestSetup {
	config := adapters.NewTestConfig(map[string]any{
		"header": "X-Request-Id",
	})

	return &requestIdMiddlewareTestSetup{
		config:            config,
		requestIdProvider: providers.NewRequestIdProvider(),
		responseRecorder:  httptest.NewRecorder(),
	}
}

func TestRequestIdMiddleware_WhenRequestIncoming_AddsRequestIdToContext(t *testing.T) {
	setup := setupRequestIdMiddlewareTest(t)

	middleware := middleware.NewRequestIdMiddleware(
		setup.config,
		setup.requestIdProvider,
	)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		value := setup.requestIdProvider.Value(r.Context())
		assert.NotNil(t, value, "request id should be in context")
		assert.IsType(t, entities.RequestId(""), *value)
		assert.NotEmpty(t, *value, "request id should not be empty")
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
}

func TestRequestIdMiddleware_WhenMultipleRequestsIncoming_RequestIdsAreUnique(t *testing.T) {
	setup := setupRequestIdMiddlewareTest(t)

	middleware := middleware.NewRequestIdMiddleware(
		setup.config,
		setup.requestIdProvider,
	)

	requestIds := make(map[entities.RequestId]struct{})
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := setup.requestIdProvider.Value(r.Context())
		assert.NotNil(t, value, "request id should be in context")
		assert.NotEmpty(t, *value, "request id should not be empty")
		requestIds[*value] = struct{}{}
	})
	handlerToTest := middleware(nextHandler)

	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "http://testing", nil)
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}

	assert.Len(t, requestIds, 100, "request ids should be unique")
}

func TestRequestIdMiddleware_WhenRequestIdInHeader_RequestIdIsInContext(t *testing.T) {
	setup := setupRequestIdMiddlewareTest(t)

	middleware := middleware.NewRequestIdMiddleware(
		setup.config,
		setup.requestIdProvider,
	)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		value := setup.requestIdProvider.Value(r.Context())
		assert.NotNil(t, value, "request id should be in context")
		assert.Equal(t, *value, entities.RequestId("test-id"), "request id should match header value")
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	req.Header.Set(setup.config.GetString("header"), "test-id")
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
}

func TestRequestIdMiddleware_WhenNoHeaderConfig_UsesDefaultHeaderName(t *testing.T) {
	setup := setupRequestIdMiddlewareTest(t)

	middleware := middleware.NewRequestIdMiddleware(
		adapters.NewTestConfig(),
		setup.requestIdProvider,
	)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		value := setup.requestIdProvider.Value(r.Context())
		assert.NotNil(t, value, "request id should be in context")
		assert.Equal(t, *value, entities.RequestId("test-id"), "request id should be the one from header")
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	req.Header.Set("X-Request-Id", "test-id")

	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
}
