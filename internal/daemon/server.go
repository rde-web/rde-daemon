package daemon

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"rde-daemon/internal/config"
	"strings"
)

type service interface {
	RegisterRoutes()
	Run() error
	Shutdown()
}

type request interface {
	lsRequest |
		catRequest |
		touchRequest |
		overrideRequest
}
type response interface {
	lsResponse |
		catResponse |
		touchResponse |
		overrideResponse
}

func makePOSTHandlerFN[I request, O response](handle func(I) (*O, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		body, errReadBody := io.ReadAll(r.Body)
		if errReadBody != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errReadBody.Error()))
			return
		}
		var req I
		if errDecode := decode(body, &req); errDecode != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errDecode.Error()))
			return
		}
		result, errHandle := handle(req)
		if errHandle != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(errHandle.Error()))
			return
		}
		resp, errEncodeResp := encode(result)
		if errEncodeResp != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errHandle.Error()))
			return
		}
		if _, errWrite := w.Write(resp); errWrite != nil {
			log.Printf("cannot write response: %v", errWrite)
		}
	}
}

func makePath(p string) string {
	var basePath = config.Instance.ProjectPath
	var dir string = path.Join(basePath, p)
	if !strings.HasPrefix(dir, basePath) {
		dir = basePath
	}
	return dir
}

func decode(data []byte, dst interface{}) error {
	// return msgpack.Unmarshal(data, dst)
	// @todo cannot correctrly decode msg pack on rde-commutator
	return json.Unmarshal(data, dst)
}

func encode(src interface{}) ([]byte, error) {
	// return msgpack.Marshal(src)
	// @todo cannot correctrly decode msg pack on rde-commutator
	return json.Marshal(src)
}

func RunService(srvc service) error {
	srvc.RegisterRoutes()
	defer srvc.Shutdown()
	return srvc.Run()
}
