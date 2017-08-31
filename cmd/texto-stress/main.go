package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kureuil/texto"
	"github.com/sirupsen/logrus"
	"github.com/satori/go.uuid"
)

func main() {
	addr := flag.String("addr", "ws://localhost:8398/v1/texto", "The address of the messaging server")
	connections := flag.Int("connections", 10000, "The number of connections to spawn")
	flag.Parse()
	log := logrus.New()
	for i := 0; i < *connections; i++ {
		go Stress(log, *addr)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<- stop
}

func Stress(log *logrus.Logger, addr string) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, make(http.Header))
	if err != nil {
		log.Fatal(err)
	}
	connectMsg := new(texto.ChatMessage)
	if err := conn.ReadJSON(connectMsg); err != nil {
		log.Fatal(err)
	}
	sendMsg := texto.ChatMessage{
		ID:       uuid.NewV4(),
		ClientID: connectMsg.ClientID,
		Kind:     texto.SendMessageKind,
		Data: texto.SendMessagePayload{
			ReceiverID: uuid.NewV4(),
			Text:       "Hello World!",
		},
	}
	marshaled, err := json.Marshal(&sendMsg)
	if err != nil {
		log.Fatal(err)
	}
	message, err := websocket.NewPreparedMessage(websocket.TextMessage, marshaled)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn.WritePreparedMessage(message)
		time.Sleep(1 * time.Second)
	}
}
