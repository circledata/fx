package fx

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type JSONAPIModule struct {
	Mux    *mux.Router
	ContextPath string
	Logger Logger
}

func (m *JSONAPIModule) SetContextPath(contextPath string){
	m.ContextPath = contextPath
}

func (m *JSONAPIModule) GetContextPath() string{
	return m.ContextPath
}

func (m *JSONAPIModule) DecodeRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (m *JSONAPIModule) SetLogger(logger Logger){
	m.Logger = logger
}

func (m *JSONAPIModule) GetLogger() Logger {
	return m.Logger
}

func (m *JSONAPIModule) EncodeResponse(w http.ResponseWriter, status int, data interface{}) error {

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

func (m *JSONAPIModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Mux.ServeHTTP(w, r)
}
