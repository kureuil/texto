package texto

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// A Broker is responsible for transmitting messages between users.
type Broker interface {
	// Regsiter adds a Client to the Broker
	Register(client *Client) error
	// Unregister removes a client from the Broker
	Unregister(client *Client) error
	// Send sends the given message to the Client associated to the given ID.
	Send(receiverID uuid.UUID, message *BrokerMessage) error
	// Poll reads all incoming messages and transmit them to known Clients.
	Poll(ctx context.Context) error
}

// A BrokerMessage is sent between Brokers to transmit the messages to the right user.
type BrokerMessage struct {
	SenderID    uuid.UUID
	RecipientID uuid.UUID
	Text        string
}

// RedisBrokerPrefix is the prefix used for all keys registered by the RedisBroker.
const RedisBrokerPrefix = "texto:"

// A RedisBroker transmits messages between users using Redis as its backend.
type RedisBroker struct {
	log     *logrus.Logger
	clients sync.Map
	conn    redis.Conn
	pubsub  redis.PubSubConn
}

// NewRedisBroker creates a new RedisBroker instance, connecting to the Redis server using the given TCP address.
func NewRedisBroker(log *logrus.Logger, addr string) (*RedisBroker, error) {
	redisConn, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	pubsubConn, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &RedisBroker{
		log:    log,
		conn:   redisConn,
		pubsub: redis.PubSubConn{Conn: pubsubConn},
	}, nil
}

// Register registers a Client in the internal Client map of the Broker.
func (b *RedisBroker) Register(client *Client) error {
	b.clients.Store(client.ID.String(), client)
	return nil
}

// Unregister removes a Client from the internal Client map of the Broker.
func (b *RedisBroker) Unregister(client *Client) error {
	b.clients.Delete(client.ID.String())
	return nil
}

// PumpMessages subscribe to Redis channels, reads all incoming messages and sends them into the given channel.
// If an error is encountered while reading the messages, the channel is closed and the function exits.
func (b *RedisBroker) PumpMessages(channelsPattern string, out chan *redis.PMessage) {
	b.pubsub.PSubscribe(channelsPattern)
	defer func() {
		b.pubsub.PUnsubscribe(channelsPattern)
	}()
	for {
		switch n := b.pubsub.Receive().(type) {
		case redis.PMessage:
			b.log.
				WithField("channel", n.Channel).
				WithField("data", n.Data).
				Info("Received pmessage")
			out <- &n
		case error:
			b.log.Error(n)
			close(out)
			return
		}
	}
}

// Poll reads all messages published on the Redis server. If a message is intended to a known user, the Broker will send
// it into the recipient's outboundChan.
func (b *RedisBroker) Poll(ctx context.Context) error {
	channelsPattern := RedisBrokerPrefix + "*"
	inbound := make(chan *redis.PMessage, 128)
	go b.PumpMessages(channelsPattern, inbound)
	for {
		select {
		case pmessage := <-inbound:
			if pmessage == nil {
				break
			}
			message := new(BrokerMessage)
			if err := json.Unmarshal(pmessage.Data, message); err != nil {
				b.log.Error(err)
				break
			}
			v, ok := b.clients.Load(message.RecipientID.String())
			if !ok {
				break
			}
			client, ok := v.(*Client)
			if !ok {
				b.log.
					WithField("sender", message.SenderID).
					WithField("recipient", message.RecipientID).
					Warn("Value is not a valid *Client")
				break
			}
			go func() {
				client.outboundChan <- NewReceiveMessage(nil, message.RecipientID, ReceiveMessagePayload{
					SenderID: message.SenderID,
					Text:     message.Text,
				})
			}()
		case <-ctx.Done():
			return nil
		}
	}
}

// Send publishes the given message on the Redis server.
func (b *RedisBroker) Send(receiverID uuid.UUID, message *BrokerMessage) error {
	marshaled, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = b.conn.Do("PUBLISH", RedisBrokerPrefix+receiverID.String(), marshaled)
	return err
}
