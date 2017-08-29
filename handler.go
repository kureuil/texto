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
	inboundChan := make(chan ChatMessage, 32)
	outboundChan := make(chan ChatMessage, 32)
	go func () {
		for {
			var message ChatMessage
			if err := conn.ReadJSON(&message); err != nil {
				h.log.Error(err)
				if err := conn.Close(); err != nil {
					h.log.Error(err)
				}
				close(inboundChan)
				return
			}
			h.log.
				WithField("kind", message.Kind).
				WithField("remote", conn.RemoteAddr()).
				Info("Received message")
			inboundChan <- message
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
				WithField("kind", outbound.Kind).
				Info("Sending message")
			if err := conn.WriteJSON(outbound); err != nil {
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
