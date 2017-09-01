package texto

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

// ChatHandler is the HTTP Handler responsible for upgrading connections to the WebSocket Protocol and managing them.
type ChatHandler struct {
	Log      *logrus.Logger
	Broker   Broker
	Upgrader websocket.Upgrader
	Timeout  time.Duration
}

// ServeHTTP is the http.Handler implementation for ChatHandler.
func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Log.Error(err)
		return
	}
	client := NewClient(h.Log, conn, h.Broker)
	h.Broker.Register(client)
	defer func() {
		h.Broker.Unregister(client)
	}()
	if len(r.URL.Query().Get("nogreet")) == 0 {
		client.outboundChan <- NewConnectionMessage(nil, client.ID, ConnectionMessagePayload{
			ClientID: client.ID,
		})
	}
	client.Run(h.Timeout)
}
