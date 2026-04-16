package router

import (
	"net/http"
	"strings"

	mdl "github.com/jmticonap/iter/middleware"
)

type GenericHandler func() (any, error)

type RequestEvent struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    map[string]any
}

type Route struct {
	Method  string
	Path    string
	Handler mdl.MdlFunc
}

type RouteSchema struct {
	Handler  mdl.MdlFunc
	Children map[string]*RouteSchema
}

type Routes struct {
	Routes map[string]*RouteSchema
}

type IRoutes interface {
	Connect(route Route) *Routes
	Delete(route Route) *Routes
	Get(route Route) *Routes
	Head(route Route) *Routes
	Options(route Route) *Routes
	Patch(route Route) *Routes
	Post(route Route) *Routes
	Put(route Route) *Routes
	Trace(route Route) *Routes
}

var rootNodeKey string = "root"

func GetRouteList() map[string]*RouteSchema {
	routes := make(map[string]*RouteSchema)
	routes[rootNodeKey] = &RouteSchema{
		Handler:  nil,
		Children: make(map[string]*RouteSchema),
	}
	methodsList := []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	}
	for _, method := range methodsList {
		rSch := &RouteSchema{
			Handler:  nil,
			Children: make(map[string]*RouteSchema),
		}
		routes[rootNodeKey].Children[method] = rSch
	}

	return routes
}

func NewRoutes() *Routes {
	return &Routes{
		Routes: GetRouteList(),
	}
}

func (r *Routes) Connect(route Route) *Routes {
	AddRoute(r.Routes[rootNodeKey].Children[http.MethodConnect], route)
	return r
}
func (r *Routes) Delete(route Route) *Routes {
	AddRoute(r.Routes[rootNodeKey].Children[http.MethodDelete], route)
	return r
}
func (r *Routes) Get(route Route) *Routes {
	AddRoute(r.Routes[rootNodeKey].Children[http.MethodGet], route)
	return r
}
func (r *Routes) Head(route Route) *Routes {
	AddRoute(r.Routes[rootNodeKey].Children[http.MethodHead], route)
	return r
}
func (r *Routes) Options(route Route) *Routes {
	AddRoute(r.Routes[rootNodeKey].Children[http.MethodOptions], route)
	return r
}
func (r *Routes) Patch(route Route) *Routes {
	AddRoute(r.Routes[rootNodeKey].Children[http.MethodPatch], route)
	return r
}
func (r *Routes) Post(route Route) *Routes {
	AddRoute(r.Routes[rootNodeKey].Children[http.MethodPost], route)
	return r
}
func (r *Routes) Put(route Route) *Routes {
	AddRoute(r.Routes[rootNodeKey].Children[http.MethodPut], route)
	return r
}
func (r *Routes) Trace(route Route) *Routes {
	AddRoute(r.Routes[rootNodeKey].Children[http.MethodTrace], route)
	return r
}

func AddRoute(rootRoute *RouteSchema, route Route) {
	cleanPath, _ := strings.CutPrefix(route.Path, "/")
	pathString := strings.Split(cleanPath, "/")
	pathLen := len(pathString)

	node, ok := rootRoute.Children[pathString[0]]
	if !ok {
		node = &RouteSchema{
			Children: make(map[string]*RouteSchema),
		}
		rootRoute.Children[pathString[0]] = node
	}

	if pathLen == 1 {
		node.Handler = route.Handler
	} else {
		AddRoute(node, Route{
			Handler: route.Handler,
			Path:    strings.Join(pathString[1:], "/"),
		})
	}
}
