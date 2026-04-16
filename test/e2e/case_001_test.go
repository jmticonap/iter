package e2e_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mdl "github.com/jmticonap/iter/middleware"
	rt "github.com/jmticonap/iter/router"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_E2E(t *testing.T) {
	req := httptest.NewRequest("POST", "http://example.com/users", nil)
	rec := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	mdlw := mdl.NewMiddleware()
	handler := func(r *http.Request) (*mdl.MidResponse, error) {
		resBody := map[string]string{"name": "Gopher", "email": "go@test.com"}
		return &mdl.MidResponse{
			Status: http.StatusCreated,
			Data:   resBody,
		}, nil
	}

	routes := rt.NewRoutes().
		Post(rt.Route{
			Path:    "/users",
			Handler: mdlw.Build(handler),
		})

	rt.HttpRouterHandler(routes.Routes)(rec, req)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, "Gopher", response["name"])
	assert.Equal(t, "go@test.com", response["email"])
	assert.Equal(t, http.StatusCreated, rec.Code)
}
