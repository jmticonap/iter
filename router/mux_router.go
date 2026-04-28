package router

import (
	"fmt"
	"net/http"

	"github.com/jmticonap/iter/middleware"
)

type MuxRouter struct {
	Mux *http.ServeMux
}

func NewMuxRouter(mux *http.ServeMux) *MuxRouter {
	var muxSrv *http.ServeMux
	if mux != nil {
		muxSrv = mux
	} else {
		muxSrv = http.NewServeMux()
	}
	return &MuxRouter{Mux: muxSrv}
}

func AdaptHandler(fn middleware.MdlFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &middleware.ResponseWrapper{ResponseWriter: w}

		_, err := fn(res, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (r *MuxRouter) Connect(route Route) *MuxRouter {
	r.Mux.HandleFunc(fmt.Sprintf("CONNECT %s", route.Path), AdaptHandler(route.Handler))
	return r
}

func (r *MuxRouter) Delete(route Route) *MuxRouter {
	r.Mux.HandleFunc(fmt.Sprintf("DELETE %s", route.Path), AdaptHandler(route.Handler))
	return r
}

func (r *MuxRouter) Get(route Route) *MuxRouter {
	r.Mux.HandleFunc(fmt.Sprintf("GET %s", route.Path), AdaptHandler(route.Handler))
	return r
}

func (r *MuxRouter) Head(route Route) *MuxRouter {
	r.Mux.HandleFunc(fmt.Sprintf("HEAD %s", route.Path), AdaptHandler(route.Handler))
	return r
}

func (r *MuxRouter) Options(route Route) *MuxRouter {
	r.Mux.HandleFunc(fmt.Sprintf("OPTIONS %s", route.Path), AdaptHandler(route.Handler))
	return r
}

func (r *MuxRouter) Patch(route Route) *MuxRouter {
	r.Mux.HandleFunc(fmt.Sprintf("PATCH %s", route.Path), AdaptHandler(route.Handler))
	return r
}

func (r *MuxRouter) Post(route Route) *MuxRouter {
	r.Mux.HandleFunc(fmt.Sprintf("POST %s", route.Path), AdaptHandler(route.Handler))
	return r
}

func (r *MuxRouter) Put(route Route) *MuxRouter {
	r.Mux.HandleFunc(fmt.Sprintf("PUT %s", route.Path), AdaptHandler(route.Handler))
	return r
}

func (r *MuxRouter) Trace(route Route) *MuxRouter {
	r.Mux.HandleFunc(fmt.Sprintf("TRACE %s", route.Path), AdaptHandler(route.Handler))
	return r
}
