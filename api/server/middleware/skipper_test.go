package middleware_test

import (
	"__template__/api/server/middleware"
	"__template__/internal/adapters"
	"__template__/internal/ports"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type skipperMiddlewareTestSetup struct {
	config           ports.Config
	responseRecorder *httptest.ResponseRecorder
}

func setupSkipperMiddlewareTest(t *testing.T) *skipperMiddlewareTestSetup {
	config := adapters.NewTestConfig(map[string]any{"skip": []string{"/skip"}})
	responseRecorder := httptest.NewRecorder()

	return &skipperMiddlewareTestSetup{
		config:           config,
		responseRecorder: responseRecorder,
	}
}

func TestSkipperMiddleware_WhenRequestPathShouldBeSkipped_SkipsMiddlewareCall(t *testing.T) {
	setup := setupSkipperMiddlewareTest(t)

	middleware := middleware.NewSkipperMiddleware(
		setup.config,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Fail(t, "middleware should be skipped")
			})
		},
	)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing/skip", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
}

func TestSkipperMiddleware_WhenRequestPathShouldNotBeSkipped_CallsNextMiddleware(t *testing.T) {
	setup := setupSkipperMiddlewareTest(t)

	middlewareCalled := false
	middleware := middleware.NewSkipperMiddleware(
		setup.config,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				middlewareCalled = true
				next.ServeHTTP(w, r)
			})
		},
	)

	handlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing/noskip", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, handlerCalled, "handler should have been called")
	assert.True(t, middlewareCalled, "middleware should have been called")
}
