package dhttp

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type HttpHandler interface {
	Handler() httprouter.Handle
}

type HttpHandlerB func(http.ResponseWriter, *http.Request, httprouter.Params)