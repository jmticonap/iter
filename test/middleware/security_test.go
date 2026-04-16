package lib_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mdl "github.com/jmticonap/iter/middleware"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareAfterBypassOnHandlerError(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/audit-check", nil)
	rec := httptest.NewRecorder()

	auditLogged := false

	_, _ = mdl.NewMiddleware().
		After(func(r *http.Request) (*mdl.MidResponse, error) {
			auditLogged = true
			return nil, nil
		}).
		Build(func(r *http.Request) (*mdl.MidResponse, error) {
			return nil, assert.AnError // Simulamos un error
		})(&mdl.ResponseWrapper{ResponseWriter: rec}, req)

	// Verificamos si la auditoría se ejecutó
	// En la nueva implementación, deferFnc se llama si err != nil
	// y dentro de deferFnc se ejecutan los MiddleAfterFuncs.
	assert.True(t, auditLogged, "VULNERABILIDAD: El middleware AFTER de auditoría fue omitido porque el handler devolvió un error")
}
