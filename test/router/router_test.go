package lib_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/jmticonap/iter/middleware"
	. "github.com/jmticonap/iter/router"
	"github.com/stretchr/testify/assert"
)

func TestHttpRouterHandler(t *testing.T) {
	routes := NewRoutes().
		Get(Route{
			Handler: NewMiddleware().
				Build(func(r *http.Request) (*MidResponse, error) {
					result := map[string]string{
						"message": "Hello, World!",
					}

					return &MidResponse{
						Status: http.StatusOK,
						Data:   result,
					}, nil
				}),
			Path: "/some/path",
		})

	req := httptest.NewRequest("GET", "http://example.com/some/path", nil)
	rec := httptest.NewRecorder()
	HttpRouterHandler(routes.Routes)(rec, req)

	assert.EqualValues(t, http.StatusOK, rec.Code)
	assert.EqualValues(t, "application/json", rec.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.EqualValues(t, "Hello, World!", response["message"])
}

func TestHttpRouterHandlerNotFound(t *testing.T) {
	routes := NewRoutes().Get(Route{
		Handler: NewMiddleware().
			Build(func(r *http.Request) (*MidResponse, error) {
				result := map[string]string{
					"message": "Hello, World!",
				}

				return &MidResponse{
					Data: result,
				}, nil
			}),
		Path: "/some/path",
	})

	req := httptest.NewRequest("GET", "http://example.com/wrong/path", nil)
	rec := httptest.NewRecorder()
	HttpRouterHandler(routes.Routes)(rec, req)

	assert.EqualValues(t, http.StatusNotFound, rec.Code)
	assert.EqualValues(t, "application/json", rec.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.EqualValues(t, "Path not found", response["message"])
}

func TestHttpRouterHandlerNil(t *testing.T) {
	routes := NewRoutes().Get(Route{
		Handler: nil,
		Path:    "/some/path",
	})

	req := httptest.NewRequest("GET", "http://example.com/some/path", nil)
	rec := httptest.NewRecorder()
	HttpRouterHandler(routes.Routes)(rec, req)

	assert.EqualValues(t, http.StatusNotFound, rec.Code)
	assert.EqualValues(t, "application/json", rec.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.EqualValues(t, "Path not found", response["message"])
}

func TestHttpRouterHandlerExecBefore(t *testing.T) {
	routes := NewRoutes().Get(Route{
		Handler: NewMiddleware().
			Build(func(r *http.Request) (*MidResponse, error) {
				result := map[string]string{
					"message": "Hello, World!",
				}

				return &MidResponse{
					Status: http.StatusOK,
					Data:   result,
				}, nil
			}),
		Path: "/some/path",
	})

	req := httptest.NewRequest("GET", "http://example.com/wrong/path", nil)
	rec := httptest.NewRecorder()
	HttpRouterHandler(routes.Routes)(rec, req)

	assert.EqualValues(t, http.StatusNotFound, rec.Code)
	assert.EqualValues(t, "application/json", rec.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.EqualValues(t, "Path not found", response["message"])
}

func TestHttpRouterHandlerDuplicateSlashes(t *testing.T) {
	routes := NewRoutes().Get(Route{
		Handler: NewMiddleware().
			Build(func(r *http.Request) (*MidResponse, error) {
				result := map[string]string{"message": "Success"}

				return &MidResponse{
					Status: http.StatusOK,
					Data:   result,
				}, nil
			}),
		Path: "/some/path",
	})

	// Petición con slashes duplicados
	req := httptest.NewRequest("GET", "http://example.com//some//path", nil)
	rec := httptest.NewRecorder()
	HttpRouterHandler(routes.Routes)(rec, req)

	assert.EqualValues(t, http.StatusOK, rec.Code, "Debería normalizar los slashes y devolver 200 OK")

	var response map[string]string
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.EqualValues(t, "Success", response["message"])
}
