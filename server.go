package texto

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"context"
)

// A Server bundles an HTTP Server and all the configuration required at runtime.
type Server struct {
	Log        *logrus.Logger
	broker     Broker
	httpServer http.Server
}

// NewServer returns an initialized Server.
func NewServer(log *logrus.Logger, addr string, broker Broker) (*Server, error) {
	mux := http.NewServeMux()
	mux.Handle("/v1/texto", &ChatHandler{
		Log:    log,
		Broker: broker,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	return &Server{
		Log: log,
		broker: broker,
		httpServer: http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadTimeout:       60 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}, nil
}

// Run tells the Server to start listening for incoming HTTP connections.
func (s *Server) Run(ctx context.Context) error {
	s.Log.WithField("addr", s.httpServer.Addr).Info("Starting HTTP server")
	go s.broker.Poll(ctx)
	return s.httpServer.ListenAndServe()
}
