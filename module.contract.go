package fx

import (
	"github.com/gorilla/mux"
)

type ServerModule interface {
	SetContextPath(contextPath string)
	GetContextPath() string
	SetLogger(logger Logger)
	GetLogger() Logger
	Initialize(*mux.Router) error
}
