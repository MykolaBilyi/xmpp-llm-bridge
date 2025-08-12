package middleware_test

import (
	"__template__/api/server/middleware"
	"__template__/internal/adapters"
	adapter_mocks "__template__/internal/adapters/mocks"
	"__template__/internal/ports"
	"__template__/internal/providers"
	"net/http"
	"net/http/httptest"
	"testing"

	extra "github.com/oxyno-zeta/gomock-extra-matcher"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type recoverMiddlewareTestSetup struct {
	config           ports.Config
	logger           *adapter_mocks.MockLogger
	loggerProvider   ports.LoggerProvider
	responseRecorder *httptest.ResponseRecorder
}

func setupRecoverMiddlewareTest(t *testing.T) *recoverMiddlewareTestSetup {
	ctrl := gomock.NewController(t)

	config := adapters.NewTestConfig(map[string]any{"stackTrace": "false"})
	logger := adapter_mocks.NewMockLogger(ctrl)
	loggerProvider := providers.NewLoggerProvider(logger)

	return &recoverMiddlewareTestSetup{
		config:           config,
		logger:           logger,
		loggerProvider:   loggerProvider,
		responseRecorder: httptest.NewRecorder(),
	}
}

func TestRecoverMiddleware_WhenHandlerPanics_ReturnsInternalServerError(t *testing.T) {
	setup := setupRecoverMiddlewareTest(t)

	middleware := middleware.NewRecoverMiddleware(
		setup.config,
		setup.loggerProvider,
	)

	setup.logger.
		EXPECT().
		Error("panic recovered", ports.Fields{"error": "panic test"})

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		panic("panic test")
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
	assert.Equal(t, http.StatusInternalServerError, setup.responseRecorder.Code, "response code should be 500")
}

func TestRecoverMiddleware_WhenHandlerPanicsAndStackTraceEnabled_LogsErrorWithStackTrace(t *testing.T) {
	setup := setupRecoverMiddlewareTest(t)

	middleware := middleware.NewRecoverMiddleware(
		adapters.NewTestConfig(map[string]any{"stackTrace": "true"}),
		setup.loggerProvider,
	)

	setup.logger.
		EXPECT().
		Error("panic recovered", extra.MapMatcher().Key("error", "panic test").Key("stack", gomock.Any()))

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		panic("panic test")
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
	assert.Equal(t, http.StatusInternalServerError, setup.responseRecorder.Code, "response code should be 500")
}
