package books

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	httpUtil "plataform/pkg/http"
)

type listerHTTP string

func (h listerHTTP) Handler() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

		resp := Response{Id: "1"}

		b, err := json.Marshal(&resp)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(b)
	}
}

func NewListerHTTP() httpUtil.Handler {
	return listerHTTP("")
}
