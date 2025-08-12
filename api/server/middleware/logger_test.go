package middleware_test

import (
	"__template__/api/server/middleware"
	"__template__/internal/adapters"
	adapter_mocks "__template__/internal/adapters/mocks"
	"__template__/internal/entities"
	"__template__/internal/ports"
	"__template__/internal/providers"
	provider_mocks "__template__/internal/providers/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type loggerMiddlewareTestSetup struct {
	config            ports.Config
	logger            *adapter_mocks.MockLogger
	loggerProvider    ports.LoggerProvider
	requestIdProvider *provider_mocks.MockRequestIdProvider
	responseRecorder  *httptest.ResponseRecorder
}

func setupLoggerMiddlewareTest(t *testing.T) *loggerMiddlewareTestSetup {
	ctrl := gomock.NewController(t)

	config := adapters.NewTestConfig()
	logger := adapter_mocks.NewMockLogger(ctrl)
	loggerProvider := providers.NewLoggerProvider(logger)
	requestIdProvider := provider_mocks.NewMockRequestIdProvider(ctrl)

	return &loggerMiddlewareTestSetup{
		config:            config,
		logger:            logger,
		loggerProvider:    loggerProvider,
		requestIdProvider: requestIdProvider,
		responseRecorder:  httptest.NewRecorder(),
	}
}

func TestLoggerMiddleware_WhenRequestIncoming_AddsLoggerToContextWithExtraLogFields(t *testing.T) {
	setup := setupLoggerMiddlewareTest(t)

	middleware := middleware.NewLoggerMiddleware(
		setup.config,
		setup.loggerProvider,
		setup.requestIdProvider,
	)

	testRequestId := entities.RequestId("test-request-id")

	setup.requestIdProvider.
		EXPECT().
		Value(gomock.Any()).
		Return(&testRequestId).
		AnyTimes()

	setup.logger.
		EXPECT().
		WithFields(gomock.Any()).
		Return(setup.logger)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "the handler should have been called")
}

func TestLoggerMiddleware_WhenNoRequestIdFieldNameInConfig_LoggerUsesDefaulFieldName(t *testing.T) {
	setup := setupLoggerMiddlewareTest(t)

	middleware := middleware.NewLoggerMiddleware(
		setup.config,
		setup.loggerProvider,
		setup.requestIdProvider,
	)

	testRequestId := entities.RequestId("test-request-id")

	setup.requestIdProvider.
		EXPECT().
		Value(gomock.Any()).
		Return(&testRequestId).
		AnyTimes()

	setup.logger.
		EXPECT().
		WithFields(ports.Fields{"request_id": testRequestId}).
		Return(setup.logger)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "the handler should have been called")
}

func TestLoggerMiddleware_WhenRequestIdFieldNameInConfig_LoggerUsesConfiguredFieldName(t *testing.T) {
	setup := setupLoggerMiddlewareTest(t)

	config := adapters.NewTestConfig(map[string]any{"requestIdKey": "request-id-test"})

	middleware := middleware.NewLoggerMiddleware(
		config,
		setup.loggerProvider,
		setup.requestIdProvider,
	)

	testRequestId := entities.RequestId("test-request-id")

	setup.requestIdProvider.
		EXPECT().
		Value(gomock.Any()).
		Return(&testRequestId).
		AnyTimes()

	setup.logger.
		EXPECT().
		WithFields(ports.Fields{"request-id-test": testRequestId}).
		Return(setup.logger)

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "the handler should have been called")
}
