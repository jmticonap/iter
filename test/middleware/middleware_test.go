package lib_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mdl "github.com/jmticonap/iter/middleware"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareBeforeResponse(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/wrong/path", nil)
	rec := httptest.NewRecorder()
	isHandlerExecute := false

	res, err := mdl.NewMiddleware().
		Use(func(r *http.Request) (*mdl.MidResponse, error) {
			result := map[string]string{
				"message": "Hello, World!",
			}

			return &mdl.MidResponse{
				Status: http.StatusOK,
				Data:   result,
			}, nil
		}).
		Build(func(r *http.Request) (*mdl.MidResponse, error) {
			isHandlerExecute = true
			return nil, nil
		})(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	result, ok := res.Data.(map[string]string)

	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, res.Data, result)
	assert.False(t, isHandlerExecute)
}

func TestMiddlewareAfter(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/wrong/path", nil)
	rec := httptest.NewRecorder()
	isAfterExecute := false

	res, err := mdl.NewMiddleware().
		Use(func(r *http.Request) (*mdl.MidResponse, error) {
			result := map[string]string{
				"message": "Hello, World!",
			}

			return &mdl.MidResponse{
				Status: http.StatusOK,
				Data:   result,
			}, nil
		}).
		After(func(r *http.Request) (*mdl.MidResponse, error) {
			isAfterExecute = true
			return nil, nil
		}).
		Build(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, nil
		})(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	result, ok := res.Data.(map[string]string)

	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, res.Data, result)
	assert.True(t, isAfterExecute)
}

func TestMiddlewareError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/wrong/path", nil)
	rec := httptest.NewRecorder()
	isHandlerExecute := false

	res, err := mdl.NewMiddleware().
		Use(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, fmt.Errorf("test error")
		}).
		Error(func(r *http.Request) (*mdl.MidResponse, error) {
			return &mdl.MidResponse{
				Status: http.StatusBadRequest,
				Data:   nil,
			}, nil
		}).
		Build(func(r *http.Request) (*mdl.MidResponse, error) {
			isHandlerExecute = true
			return nil, nil
		})(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	assert.NoError(t, err)
	assert.Equal(t, res.Status, http.StatusBadRequest)
	assert.False(t, isHandlerExecute)
}

func TestMiddlewareErrorFail(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/wrong/path", nil)
	rec := httptest.NewRecorder()
	isHandlerExecute := false

	res, err := mdl.NewMiddleware().
		Use(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, fmt.Errorf("test error")
		}).
		Error(func(r *http.Request) (*mdl.MidResponse, error) {
			return &mdl.MidResponse{
				Status: http.StatusBadRequest,
				Data:   nil,
			}, nil
		}).
		Error(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, fmt.Errorf("test error fail")
		}).
		Build(func(r *http.Request) (*mdl.MidResponse, error) {
			isHandlerExecute = true
			return nil, nil
		})(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.False(t, isHandlerExecute)
}

func TestMiddlewareBeforeError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/wrong/path", nil)
	rec := httptest.NewRecorder()

	_, err := mdl.NewMiddleware().
		Use(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, fmt.Errorf("Some error")
		}).
		Build(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, nil
		})(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	assert.EqualError(t, err, "Some error")
}

func TestMiddlewareAfterError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/wrong/path", nil)
	rec := httptest.NewRecorder()

	afterExecuted := false

	_, err := mdl.NewMiddleware().
		After(func(r *http.Request) (*mdl.MidResponse, error) {
			afterExecuted = true
			return nil, fmt.Errorf("After error")
		}).
		Build(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, fmt.Errorf("Main error")
		})(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	assert.EqualError(t, err, "Main error")
	assert.True(t, afterExecuted)
}

func TestMiddlewareFlowBypassVulnerability(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/sensitive-data", nil)
	rec := httptest.NewRecorder()

	handlerExecuted := false

	_, err := mdl.NewMiddleware().
		Use(func(r *http.Request) (*mdl.MidResponse, error) {
			return &mdl.MidResponse{
				Data: "Unauthorized",
			}, fmt.Errorf("Unauthorized")
		}).
		Use(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, nil
		}).
		Build(func(r *http.Request) (*mdl.MidResponse, error) {
			handlerExecuted = true
			return &mdl.MidResponse{Data: "Sensitive Data"}, nil
		})(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	assert.Error(t, err)
	assert.False(t, handlerExecuted, "VULNERABILIDAD: El handler se ejecutó a pesar de que el primer middleware falló")
}

func TestMiddleware_ShouldStopIfBeforeReturnsResponse(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handlerExecuted := false
	afterExecuted := false

	m := mdl.NewMiddleware().
		Use(func(r *http.Request) (*mdl.MidResponse, error) {
			// Returns a response with Status and Data
			return &mdl.MidResponse{
				Status: http.StatusForbidden,
				Data:   map[string]string{"error": "stopped"},
			}, nil
		}).
		After(func(r *http.Request) (*mdl.MidResponse, error) {
			afterExecuted = true
			return nil, nil
		})
		// func(ctx context.Context, r *http.Request) (*MidResponse, error)
	handler := m.Build(func(r *http.Request) (*mdl.MidResponse, error) {
		handlerExecuted = true
		return nil, nil
	})

	_, err := handler(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	assert.NoError(t, err)
	assert.False(t, handlerExecuted, "Handler should NOT have been executed")
	assert.True(t, afterExecuted, "After middleware SHOULD have been executed")
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), "stopped")
}

func TestMiddleware_ShouldAlwaysRunAfter(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	afterExecuted := false

	m := mdl.NewMiddleware().
		Use(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, nil
		}).
		After(func(r *http.Request) (*mdl.MidResponse, error) {
			afterExecuted = true
			return nil, nil
		})

	handler := m.Build(func(r *http.Request) (*mdl.MidResponse, error) {
		return nil, nil
	})

	handler(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	assert.True(t, afterExecuted, "After middleware SHOULD have been executed")
}
