package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	mdl "github.com/jmticonap/iter/middleware"
)

func GetHandler(route *RouteSchema, path string) (mdl.MdlFunc, error) {
	cleanPath, _ := strings.CutPrefix(path, "/")
	pathString := strings.Split(cleanPath, "/")
	pathLen := len(pathString)

	var handler mdl.MdlFunc

	if rt, ok := route.Children[pathString[0]]; ok {
		if pathLen == 1 {
			handler = route.Children[pathString[0]].Handler
		} else {
			return GetHandler(
				rt, strings.Join(pathString[1:], "/"),
			)
		}
	} else {
		return nil, fmt.Errorf("Path not found\n")
	}

	return handler, nil
}

func HttpRouterHandler(routes map[string]*RouteSchema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := path.Clean(r.URL.Path)

		var handler mdl.MdlFunc
		var err error

		handler, err = GetHandler(routes[rootNodeKey].Children[method], path)

		if err == nil && handler != nil {
			handler(&mdl.ResponseWrapper{ResponseWriter: w}, r)
			return
		} else {
			msg := "Path not found"
			result := map[string]string{
				"message": msg,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(result)

			log.Println(msg)

			return
		}
	}
}
