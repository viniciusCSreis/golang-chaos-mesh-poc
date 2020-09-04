package books

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"plataform/pkg/dhttp"
)

type listerHTTP string

func (h listerHTTP) Handler() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

		resp := Response{Id: uuid.New().String()}

		b, err := json.Marshal(&resp)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}

		_, _ = writer.Write(b)
	}
}

func NewListerHTTP() dhttp.HttpHandler {
	return listerHTTP("")
}
