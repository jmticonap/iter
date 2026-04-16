package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"unsafe"
)

type ResponseWrapper struct {
	http.ResponseWriter
	written bool
}

func (res *ResponseWrapper) SetHeader(key, value string) {
	if res.written {
		return
	}
	res.ResponseWriter.Header().Set(key, value)
}

func (res *ResponseWrapper) IsWritten() bool {
	return res.written
}

func (res *ResponseWrapper) WriteResponse(status int, obj any) error {
	if res.written {
		return fmt.Errorf("Response already write")
	}

	res.ResponseWriter.WriteHeader(status)
	var err error
	val, ok := obj.(string)
	if ok {
		_, err = res.ResponseWriter.Write(StringToByte(val))
	} else {
		err = json.NewEncoder(res.ResponseWriter).Encode(obj)
	}
	res.written = true

	return err
}

func StringToByte(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s)) // do not copy
}
