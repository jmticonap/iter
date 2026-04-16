package lib_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mdl "github.com/jmticonap/iter/middleware"
	"github.com/stretchr/testify/assert"
)

func TestResponseWrapperSetHeader(t *testing.T) {
	rec := httptest.NewRecorder()

	wrapper := mdl.ResponseWrapper{
		ResponseWriter: rec,
	}

	t.Run("SetHeader", func(t *testing.T) {
		wrapper.SetHeader("test-key", "test-value")

		assert.Equal(t, "test-value", wrapper.ResponseWriter.Header().Get("test-key"))
	})

	t.Run("SetHeader after write response", func(t *testing.T) {
		err := wrapper.WriteResponse(http.StatusOK, "test-response")
		wrapper.SetHeader("test-key-2", "test-value-2")

		assert.NoError(t, err)
		assert.Equal(t, "", wrapper.ResponseWriter.Header().Get("test-key-2"))
	})

	rec2 := httptest.NewRecorder()

	wrapper2 := mdl.ResponseWrapper{
		ResponseWriter: rec2,
	}

	t.Run("SetHeader after write response blank", func(t *testing.T) {
		err := wrapper2.WriteResponse(http.StatusOK, "")
		wrapper2.SetHeader("test-key", "test-value")

		assert.True(t, wrapper2.IsWritten())
		assert.NoError(t, err)
		assert.Equal(t, "", wrapper2.ResponseWriter.Header().Get("test-key"))
	})

	t.Run("error at WriteRespon at second time", func(t *testing.T) {
		_ = wrapper2.WriteResponse(http.StatusOK, "")
		err := wrapper2.WriteResponse(http.StatusOK, "")

		assert.Error(t, err, "WriteResponse should response error")
	})
}
