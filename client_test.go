package texto

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestClient_HandleMessage(t *testing.T) {
	client := NewClient(nil, nil)

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
