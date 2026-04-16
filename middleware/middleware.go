package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type MidResponse struct {
	Status  int
	Headers *map[string]string
	Data    any
}

type MdlNotWritableFunc func(r *http.Request) (*MidResponse, error)
type MdlFunc func(res *ResponseWrapper, r *http.Request) (*MidResponse, error)

type Middleware struct {
	MiddleFuncs      []MdlNotWritableFunc
	MiddleAfterFuncs []MdlNotWritableFunc
	MiddleErrorFuncs []MdlNotWritableFunc
}

type IMiddleware interface {
	Use(MdlNotWritableFunc) *IMiddleware
	After(MdlNotWritableFunc) *IMiddleware
	Error(MdlNotWritableFunc) *IMiddleware
	Build(MdlNotWritableFunc) MdlFunc
}

func NewMiddleware() *Middleware {
	return &Middleware{
		MiddleFuncs:      []MdlNotWritableFunc{},
		MiddleAfterFuncs: []MdlNotWritableFunc{},
		MiddleErrorFuncs: []MdlNotWritableFunc{},
	}
}

func (m *Middleware) Use(cnf MdlNotWritableFunc) *Middleware {
	m.MiddleFuncs = append(m.MiddleFuncs, cnf)
	return m
}

func (m *Middleware) After(fn MdlNotWritableFunc) *Middleware {
	m.MiddleAfterFuncs = append(m.MiddleAfterFuncs, fn)
	return m
}

func (m *Middleware) Error(fn MdlNotWritableFunc) *Middleware {
	m.MiddleErrorFuncs = append(m.MiddleErrorFuncs, fn)
	return m
}

func (m *Middleware) Build(fn MdlNotWritableFunc) MdlFunc {
	return func(res *ResponseWrapper, r *http.Request) (*MidResponse, error) {
		res.SetHeader("Content-Type", "application/json")
		var result *MidResponse
		var err error

		executeAfters := func() {
			for _, middleAfterFunc := range m.MiddleAfterFuncs {
				_, aErr := middleAfterFunc(r)
				if aErr != nil {
					slog.Error("Error in after middleware", slog.String("error", aErr.Error()))
					res.WriteHeader(http.StatusInternalServerError)
					res.Write([]byte("{}"))
					break
				}
			}
		}

		handleError := func(origRes *MidResponse, origErr error) (*MidResponse, error) {
			resEr, er := origRes, origErr
			for _, middleErrorFunc := range m.MiddleErrorFuncs {
				tmpRes, tmpErr := middleErrorFunc(r)
				if tmpErr != nil {
					slog.Error("Error in error middleware", slog.String("error", tmpErr.Error()))
					er = tmpErr
					resEr = nil
					break
				}
				if tmpRes != nil && tmpRes.Status > 0 {
					res.WriteHeader(tmpRes.Status)
					if err := json.NewEncoder(res.ResponseWriter).Encode(tmpRes.Data); err != nil {
						jsonError := fmt.Sprintf(`{"error": "%s"}`, err.Error())
						res.Write([]byte(jsonError))
					}
					resEr, er = tmpRes, tmpErr
				}
			}
			executeAfters()
			return resEr, er
		}

		// Before flow middleware
		for _, middleFunc := range m.MiddleFuncs {
			result, err = middleFunc(r)

			if err != nil {
				slog.Error("Error in middleware", slog.String("error", err.Error()))
				return handleError(result, err)
			}
			if result != nil && result.Data != nil && result.Status > 0 {
				res.WriteResponse(result.Status, result.Data)
				executeAfters()
				return result, err
			}
		}

		// Handler
		result, err = fn(r)
		if !res.written && (result != nil && result.Status == 0) {
			if err != nil {
				slog.Error("Error in Handler:", slog.String("error", err.Error()))
			} else {
				slog.Error("Error in Handler:")
			}

			res.WriteResponse(http.StatusInternalServerError, nil)
		} else if !res.written && (result != nil && result.Status > 0) {
			res.WriteResponse(result.Status, result.Data)
		}

		if err != nil {
			return handleError(result, err)
		}

		executeAfters()
		return result, err
	}
}
