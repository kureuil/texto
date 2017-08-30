package texto

import (
	"net/http"

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
	client := NewClient(conn, h.log)
	client.Run()
}
