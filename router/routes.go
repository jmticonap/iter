package router

import (
	mdl "github.com/jmticonap/iter/middleware"
)

type Route struct {
	Method  string
	Path    string
	Handler mdl.MdlFunc
}
