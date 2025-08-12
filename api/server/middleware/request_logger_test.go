package middleware_test

import (
	"__template__/api/server/middleware"
	"__template__/internal/adapters"
	adapter_mocks "__template__/internal/adapters/mocks"
	"__template__/internal/ports"
	"__template__/internal/providers"
	provider_mocks "__template__/internal/providers/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type requestLoggerMiddlewareTestSetup struct {
	config            ports.Config
	logger            *adapter_mocks.MockLogger
	loggerProvider    ports.LoggerProvider
	requestIdProvider *provider_mocks.MockRequestIdProvider
	responseRecorder  *httptest.ResponseRecorder
}

func setupRequestLoggerMiddlewareTest(t *testing.T) *requestLoggerMiddlewareTestSetup {
	ctrl := gomock.NewController(t)

	config := adapters.NewTestConfig(map[string]any{"level": "info"})
	logger := adapter_mocks.NewMockLogger(ctrl)
	loggerProvider := providers.NewLoggerProvider(logger)
	requestIdProvider := provider_mocks.NewMockRequestIdProvider(ctrl)

	return &requestLoggerMiddlewareTestSetup{
		config:            config,
		logger:            logger,
		loggerProvider:    loggerProvider,
		requestIdProvider: requestIdProvider,
		responseRecorder:  httptest.NewRecorder(),
	}
}

func TestRequestLoggerMiddleware_WhenRequestIncoming_LogsIncomingRequest(t *testing.T) {
	setup := setupRequestLoggerMiddlewareTest(t)

	middleware := middleware.NewRequestLoggerMiddleware(
		setup.config,
		setup.loggerProvider,
	)

	setup.logger.
		EXPECT().
		Info("request", ports.Fields{"method": "GET", "path": "/test"}).
		Times(1)

	setup.logger.
		EXPECT().
		Info("response", gomock.Any()).
		Times(1)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing/test", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
}

func TestRequestLoggerMiddleware_WhenNoLevelInConfig_LogsLevelDebug(t *testing.T) {
	setup := setupRequestLoggerMiddlewareTest(t)

	middleware := middleware.NewRequestLoggerMiddleware(
		adapters.NewTestConfig(),
		setup.loggerProvider,
	)

	setup.logger.
		EXPECT().
		Debug("request", ports.Fields{"method": "GET", "path": "/test"}).
		Times(1)

	setup.logger.
		EXPECT().
		Debug("response", gomock.Any()).
		Times(1)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing/test", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
}

func TestRequestLoggerMiddleware_WhenUnknownLevelInConfig_LogsLevelInfo(t *testing.T) {
	setup := setupRequestLoggerMiddlewareTest(t)

	setup.logger.
		EXPECT().
		Warn("invalid log level for request logging, using info", ports.Fields{"value": "unknown"}).
		Times(1)

	middleware := middleware.NewRequestLoggerMiddleware(
		adapters.NewTestConfig(map[string]any{"level": "unknown"}),
		setup.loggerProvider,
	)

	setup.logger.
		EXPECT().
		Info("request", ports.Fields{"method": "GET", "path": "/test"}).
		Times(1)

	setup.logger.
		EXPECT().
		Info("response", gomock.Any()).
		Times(1)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing/test", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
}
