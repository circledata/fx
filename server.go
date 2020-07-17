package fx

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(options ...ServerOption) (*BaseTwoServer, error) {
	server := &BaseTwoServer{
		mux:           mux.NewRouter().StrictSlash(true),
		serverAddress: ":8888",
	}

	if len(options) > 0 {
		for _, option := range options {
			option(server)
		}
	}

	validationError := server.validateOptions()

	if validationError != nil {
		return nil, validationError
	}

	return server, nil
}

type ServerOption func(*BaseTwoServer)

func SetAddress(addr string) ServerOption {
	return func(server *BaseTwoServer) {
		server.serverAddress = addr
	}
}

func SetLogger(logger Logger) ServerOption {
	return func(server *BaseTwoServer) {
		server.logger = logger
	}
}

func SetReadTimeout(readTimeout time.Duration) ServerOption {
	return func(server *BaseTwoServer) {
		server.readTimeout = readTimeout
	}
}

func SetWriteTimeout(writeTimeout time.Duration) ServerOption {
	return func(server *BaseTwoServer) {
		server.writeTimeout = writeTimeout
	}
}

func SetIdleTimeout(idleTimeout time.Duration) ServerOption {
	return func(server *BaseTwoServer) {
		server.idleTimeout = idleTimeout
	}
}

func SetTLSConfig(tlsConfig *tls.Config) ServerOption {
	return func(server *BaseTwoServer) {
		server.tlsConfig = tlsConfig
	}
}

type BaseTwoServer struct {
	mux           *mux.Router
	serverAddress string
	logger        Logger
	readTimeout   time.Duration
	writeTimeout  time.Duration
	idleTimeout   time.Duration
	tlsConfig     *tls.Config
}

func (s *BaseTwoServer) validateOptions() error {

	if s.logger == nil {
		return errors.New("server logger has not been set")
	}

	return nil
}

func (s *BaseTwoServer) GetMux() *mux.Router {
	return s.mux
}

func (s *BaseTwoServer) RegisterModule(contextPath string, module ServerModule) error {
	router := s.mux.PathPrefix(contextPath).Subrouter().StrictSlash(true)

	module.SetContextPath(contextPath)
	module.SetLogger(s.logger)

	moduleInitError := module.Initialize(router)

	if moduleInitError != nil {
		return moduleInitError
	}

	return nil
}

func (s *BaseTwoServer) HandlePanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("unknown error")
			}

			s.logger.Error(err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
		}()
		h.ServeHTTP(w, r)
	})
}

func (s *BaseTwoServer) Run() {
	srv := &http.Server{
		Addr:         s.serverAddress,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
		TLSConfig:    s.tlsConfig,
		Handler:      http.TimeoutHandler(s.HandlePanic(s.mux), 1*time.Minute, "Service unavailable. Request timeout"),
	}

	log.Println("Server is now running on " + s.serverAddress)
	log.Println(srv.ListenAndServe())
}
