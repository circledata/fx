package fx

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type APIModule struct {
	Mux    *mux.Router
}

func (s *APIModule) DecodeRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *APIModule) EncodeResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if data == nil {
		return nil
	}

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		return err
	}

	return nil
}

func (m *APIModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Mux.ServeHTTP(w, r)
}
