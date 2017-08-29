package texto

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// ChatHandler is the HTTP Handler responsible for upgrading connections to the WebSocket Protocol and managing them.
type ChatHandler struct {
	log      *logrus.Logger
	upgrader websocket.Upgrader
}

// ServeHTTP is the http.Handler implementation for ChatHandler.
func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error(err)
		return
	}
	inboundChan := make(chan []byte, 64)
	outboundChan := make(chan []byte, 64)
	go func () {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				h.log.Error(err)
				if err := conn.Close(); err != nil {
					h.log.Error(err)
				}
				close(inboundChan)
				return
			}
			h.log.
				WithField("type", messageType).
				WithField("length", len(p)).
				WithField("remote", conn.RemoteAddr()).
				Info("Received message")
			inboundChan <- p
		}
	}()
	for {
		select {
		case inbound := <-inboundChan:
			// Just an echo server for now
			go func() {
				outboundChan <- inbound
			}()
		case outbound := <-outboundChan:
			h.log.
				WithField("remote", conn.RemoteAddr()).
				WithField("length", len(outbound)).
				Info("Sending message")
			if err := conn.WriteMessage(websocket.TextMessage, outbound); err != nil {
				h.log.Error(err)
				return
			}
		case <-time.After(5 * time.Minute):
			h.log.Info("Connection timeout")
			if err := conn.Close(); err != nil {
				h.log.Error(err)
			}
			break
		}
	}
}

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
