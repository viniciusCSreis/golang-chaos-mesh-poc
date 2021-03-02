package http

import (
	"github.com/julienschmidt/httprouter"
)

const AcceptHeader = "Accept"

type Handler interface {
	Handler() httprouter.Handle
}
