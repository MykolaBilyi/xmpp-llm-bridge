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
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type timeoutMiddlewareTestSetup struct {
	config           ports.Config
	logger           *adapter_mocks.MockLogger
	loggerProvider   ports.LoggerProvider
	responseRecorder *httptest.ResponseRecorder
}

func setupTimeoutMiddlewareTest(t *testing.T) *timeoutMiddlewareTestSetup {
	ctrl := gomock.NewController(t)

	config := adapters.NewTestConfig(map[string]any{"timeout": "100ms"})
	logger := adapter_mocks.NewMockLogger(ctrl)
	loggerProvider := providers.NewLoggerProvider(logger)

	return &timeoutMiddlewareTestSetup{
		config:           config,
		logger:           logger,
		loggerProvider:   loggerProvider,
		responseRecorder: httptest.NewRecorder(),
	}
}

func TestTimeoutMiddleware_WhenTimeoutReaches_ReturnsRequestTimeout(t *testing.T) {
	setup := setupTimeoutMiddlewareTest(t)

	middleware := middleware.NewTimeoutMiddleware(
		setup.config,
		setup.loggerProvider,
	)

	setup.logger.
		EXPECT().
		Warn("request timeout")

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		time.Sleep(200 * time.Millisecond)
	})
	handlerToTest := middleware(nextHandler)

	req := httptest.NewRequest("GET", "http://testing", nil)
	handlerToTest.ServeHTTP(setup.responseRecorder, req)

	assert.True(t, called, "handler should have been called")
	assert.Equal(t, http.StatusRequestTimeout, setup.responseRecorder.Code, "response code should be 408")
}
