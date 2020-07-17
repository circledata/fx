package fx

import (
	"encoding/xml"
	"net/http"

	"github.com/gorilla/mux"
)

type XMLAPIModule struct {
	Mux    *mux.Router
	ContextPath string
	Logger Logger
}

func (m *XMLAPIModule) SetContextPath(contextPath string){
	m.ContextPath = contextPath
}

func (m *XMLAPIModule) GetContextPath() string{
	return m.ContextPath
}

func (m *XMLAPIModule) DecodeRequest(r *http.Request, v interface{}) error {
	return xml.NewDecoder(r.Body).Decode(v)
}

func (m *XMLAPIModule) SetLogger(logger Logger){
	m.Logger = logger
}

func (m *XMLAPIModule) GetLogger() Logger {
	return m.Logger
}

func (m *XMLAPIModule) EncodeResponse(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)

	if data == nil {
		return nil
	}

	err := xml.NewEncoder(w).Encode(data)

	if err != nil {
		return err
	}

	return nil
}

func (m *XMLAPIModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Mux.ServeHTTP(w, r)
}
