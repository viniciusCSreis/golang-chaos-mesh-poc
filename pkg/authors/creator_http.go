package authors

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
	"plataform/pkg/api"
	"plataform/pkg/dhttp"
	"plataform/pkg/provider/messaging"
)

type creatorHTTP struct {
	m messaging.PublisherSync
}

func (h creatorHTTP) Handler() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

		resp := Response{Id: uuid.New().String()}

		b, err := json.Marshal(&resp)
		if err != nil {
			log.Error().Err(err)
			writer.WriteHeader(http.StatusInternalServerError)
		}

		msg := messaging.Message{
			Content: b,
		}

		err = h.m.PublishSync(api.DefaultOrganization, messaging.SubjectBuildBooK, msg)
		if err != nil {
			log.Error().Err(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(b)
	}
}

func NewCreatorHTTP(msg messaging.PublisherSync) dhttp.HttpHandler {
	return creatorHTTP{msg}
}
