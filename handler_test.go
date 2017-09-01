package texto

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func newLogger() *logrus.Logger {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log
}

func TestChatHandler_ServeHTTP(t *testing.T) {
	handler := ChatHandler{
		Log: newLogger(),
		Broker: newDummyBroker(),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Timeout: 3 * time.Second,
	}
	srv := httptest.NewServer(&handler)
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"
	u.Path = "/v1/texto"
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.Nil(t, err)
	greetingMsg := new(ChatMessage)
	err = conn.ReadJSON(greetingMsg)
	assert.Nil(t, err)
	assert.Equal(t, greetingMsg.Kind, ConnectionMessageKind)
	time.Sleep(3 * time.Second)
	srv.Close()
}
