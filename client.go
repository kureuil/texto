package texto

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// A Client represents an open WebSocket connection with a user.
type Client struct {
	// The universally unique ID of the user on this node.
	ID uuid.UUID

	// The current logrus instance.
	log *logrus.Logger

	// The WebSocket connection associated to this user.
	conn *websocket.Conn

	// The Broker in which the client is registered
	broker Broker

	// inboundChan is used to transfer messages incoming from the user to the main client loop.
	inboundChan chan *ChatMessage

	// outboundChan is used to transfer messages incoming from the broker to the main client loop.
	outboundChan chan *ChatMessage
}

// NewClient creates a new Client from an open WebSocket connection.
func NewClient(log *logrus.Logger, conn *websocket.Conn, broker Broker) *Client {
	return &Client{
		ID:           uuid.NewV4(),
		log:          log,
		broker:       broker,
		conn:         conn,
		inboundChan:  make(chan *ChatMessage, 32),
		outboundChan: make(chan *ChatMessage, 32),
	}
}

// consumeWebsocket reads incoming messages from the socket, and transfer them to the main client loop using the
// inboundChan channel.
func (c *Client) consumeWebsocket() {
	for {
		message := new(ChatMessage)
		if err := c.conn.ReadJSON(message); err != nil {
			c.log.Error(err)
			if _, ok := err.(*websocket.CloseError); ok {
				if err := c.conn.Close(); err != nil {
					c.log.Error(err)
				}
				close(c.inboundChan)
				break
			}
			c.outboundChan <- NewErrorMessage(nil, c.ID, ErrorMessagePayload{
				Code:        "ESYNTAX",
				Description: "Unable to process the message due to a syntax error.",
			})
			continue
		}
		c.inboundChan <- message
	}
}

// HandleMessage processes the given message and returns the ChatMessage that should be send back to the user.
func (c *Client) HandleMessage(msg *ChatMessage) *ChatMessage {
	switch msg.Kind {
	case ErrorMessageKind: // Ignore incoming error messages
	case AcknowledgeMessageKind: // Ignore incoming ack messages
	case RegistrationKind:
		return NewConnectionMessage(&msg.ID, c.ID, ConnectionMessagePayload{
			ClientID: c.ID,
		})
	case SendMessageKind:
		if msg.ClientID != c.ID {
			return NewErrorMessage(&msg.ID, c.ID, ErrorMessagePayload{
				Code:        "ECID",
				Description: "The submitted client ID doesn't match the current session.",
			})
		}
		payload, ok := msg.Data.(SendMessagePayload)
		if !ok {
			return NewErrorMessage(&msg.ID, c.ID, ErrorMessagePayload{
				Code: "EINVAL",
				Description: "The data payload doesn't match the given kind",
			})
		}
		if err := c.broker.Send(payload.ReceiverID, &BrokerMessage{
			SenderID: c.ID,
			RecipientID: payload.ReceiverID,
			Text: payload.Text,
		}); err != nil {
			c.log.Error(err)
			return NewErrorMessage(&msg.ID, c.ID, ErrorMessagePayload{
				Code: "EBROKER",
				Description: "Unable to send the message to the recipient.",
			})
		}
		return NewAckMessage(&msg.ID, c.ID)
	default:
		return NewErrorMessage(&msg.ID, c.ID, ErrorMessagePayload{
			Code:        "EKIND",
			Description: "Invalid message kind received.",
		})
	}
	return nil
}

// Run listens on the inboundChan and outboundChan for new messages to process or send.
// It timeouts after 5 minutes of inactivity.
func (c *Client) Run(timeout time.Duration) {
	go c.consumeWebsocket()
	for {
		select {
		case inbound := <-c.inboundChan:
			if inbound == nil {
				return
			}
			c.log.
				WithField("client", c.ID.String()).
				WithField("kind", inbound.Kind).
				WithField("remote", c.conn.RemoteAddr()).
				Info("Received message")
			if response := c.HandleMessage(inbound); response != nil {
				go func() {
					c.outboundChan <- response
				}()
			}
		case outbound := <-c.outboundChan:
			c.log.
				WithField("client", c.ID.String()).
				WithField("remote", c.conn.RemoteAddr()).
				WithField("kind", outbound.Kind).
				Info("Sending message")
			if err := c.conn.WriteJSON(outbound); err != nil {
				c.log.Error(err)
				return
			}
		case <-time.After(timeout):
			c.log.
				WithField("client", c.ID.String()).
				Info("Connection timeout")
			if err := c.conn.Close(); err != nil {
				c.log.Error(err)
			}
			return
		}
	}
}
