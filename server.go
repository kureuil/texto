//go:generate statik -f -src=./public

package texto

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	_ "github.com/kureuil/texto/statik"
)

// A Server bundles an HTTP Server and all the configuration required at runtime.
type Server struct {
	Log        *logrus.Logger
	Broker     Broker
	HTTPServer http.Server
	cancelFunc context.CancelFunc
	ctx        context.Context
}

// NewServer returns an initialized Server.
func NewServer(parent context.Context, log *logrus.Logger, addr string, broker Broker) (*Server, error) {
	ctx, cancel := context.WithCancel(parent)
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
		Timeout: 5 * time.Minute,
	})
	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}
	mux.Handle("/", http.FileServer(statikFS))
	return &Server{
		Log:    log,
		Broker: broker,
		HTTPServer: http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadTimeout:       60 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		cancelFunc: cancel,
		ctx:        ctx,
	}, nil
}

// Run tells the Server to start listening for incoming HTTP connections.
func (s *Server) Run() error {
	s.Log.WithField("addr", s.HTTPServer.Addr).Info("Starting HTTP server")
	go s.Broker.Poll(s.ctx)
	return s.HTTPServer.ListenAndServe()
}

// Stop gracefully stops the server.
func (s *Server) Stop() error {
	if err := s.HTTPServer.Shutdown(s.ctx); err != nil {
		return err
	}
	s.cancelFunc()
	return nil
}
