package texto

import (
	"context"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type DummyBroker struct {}

func (b *DummyBroker) Register(client *Client) error {
	return nil
}

func (b *DummyBroker) Unregister(client *Client) error {
	return nil
}

func (b *DummyBroker) Send(receiverID uuid.UUID, message *BrokerMessage) error {
	return nil
}

func (b *DummyBroker) Poll(ctx context.Context) error {
	return nil
}

func TestClient_HandleMessage(t *testing.T) {
	client := NewClient(nil, nil, &DummyBroker{})

	errorMsg := NewErrorMessage(nil, client.ID, ErrorMessagePayload{
		Code: "ENOMEM",
		Description: "Out-of-memory",
	})
	assert.Nil(t, client.HandleMessage(errorMsg))

	ackMsg := NewAckMessage(nil, client.ID)
	assert.Nil(t, client.HandleMessage(ackMsg))

	registrationMsg := NewRegistrationMessage(nil, client.ID)
	registrationAnswer := client.HandleMessage(registrationMsg)
	assert.Equal(t, registrationMsg.ID, registrationAnswer.ID)
	assert.Equal(t, ConnectionMessageKind, registrationAnswer.Kind)

	sendMsg := NewSendMessage(nil, client.ID, SendMessagePayload{
		ReceiverID: client.ID,
		Text: "Hello World!",
	})
	sendAnswer := client.HandleMessage(sendMsg)
	assert.Equal(t, sendMsg.ID, sendAnswer.ID)
	assert.Equal(t, AcknowledgeMessageKind, sendAnswer.Kind)
}
