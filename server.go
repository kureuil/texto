package texto

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// A Server bundles an HTTP Server and all the configuration required at runtime.
type Server struct {
	log        *logrus.Logger
	httpServer http.Server
}

// NewServer returns an initialized Server.
func NewServer(addr string, logger *logrus.Logger) Server {
	mux := http.NewServeMux()
	mux.Handle("/v1/texto", &ChatHandler{
		log: logger,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	return Server{
		log: logger,
		httpServer: http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadTimeout:       60 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

// Run tells the Server to start listening for incoming HTTP connections.
func (s Server) Run() error {
	s.log.WithField("addr", s.httpServer.Addr).Info("Starting HTTP server")
	return s.httpServer.ListenAndServe()
}
