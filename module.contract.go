package fx

import (
	"github.com/gorilla/mux"
)

type ServerModule interface {
	SetContextPath(contextPath string)
	Initialize(*mux.Router) error
}
